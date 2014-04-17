#!/usr/bin/env python
#
# Copyright 2011 Hunter Freyer and Michael Kelly
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

import collections
import datetime
import hmac
import json
import logging
import urllib2

from google.appengine.api import urlfetch
from google.appengine.api import users
from google.appengine.ext import db

import models
import config_NOCOMMIT as config

# Whether to display the new vote counts (cached in Triples).
USE_CACHED_VOTE_COUNTS = True

Image = collections.namedtuple('Image', ['original', 'thumbnail'])

class EntityValidationError(Exception): pass


def GetStatsUrlsForTriple(triple, w=160, h=85):
  """Returns a list of stats URLs for the given triple.

  Args:
    triple: (Triple) triple to examine
    w: (int) Optional. Width of each chart image.
    h: (int) Optional. Height of each chart image.

  Returns:
    [str, str, str]: URLs for the Triple's three Entities.
  """
  counts = GetTripleVoteCounts(triple)

  urls = []
  overall_max = max([max(c) for c in counts])
  for count in counts:
    urls.append('http://chart.apis.google.com/chart'
                '?chxr=0,0,%(max)d'
                '&chxt=y'
                '&chbh=a'
                '&chs=%(w)dx%(h)d'
                '&cht=bvg'
                '&chco=9911BB,C76FDD,63067A'
                '&chds=0,%(max)d,0,%(max)d,0,%(max)d'
                '&chd=t:%(m)d|%(f)d|%(k)d'
                '&chdl=Marry|Fuck|Kill'
                '&chdlp=r' % (dict(m=count[0], f=count[1], k=count[2],
                              max=overall_max,
                              w=w,
                              h=h)))
  return urls

def GetTripleVoteCounts(triple):
  """Calculates vote count for the given triple.

  Returns:
    ([[int]]) Vote counts. This is a nested list (first level is votes for
        entity one, two, three; second level is votes for m, f, k).
  """
  def _CalculateEntityVoteCounts(entity):
    m = entity.assignment_reference_marry_set.count()
    f = entity.assignment_reference_fuck_set.count()
    k = entity.assignment_reference_kill_set.count()
    return [m, f, k]

  # For backwards compatibility with Triples that don't have embedded vote
  # counts.
  if not USE_CACHED_VOTE_COUNTS or not triple.has_cached_votes:
    logging.info('Updating legacy Triple without vote counts: %s',
                 triple.key())
    votes = [_CalculateEntityVoteCounts(triple.one),
             _CalculateEntityVoteCounts(triple.two),
             _CalculateEntityVoteCounts(triple.three)]
    # Race condition here: We're done calculating the votes, and we're about to
    # update the Entity. We might be off by one if someone else votes while
    # we're here. We have MapReduces to fix this up, so we don't care too much.
    db.run_in_transaction(_UpdateTripleVoteCounts, triple.key(), votes)
    return votes
  else:
    logging.info('Got cached votes for Triple %s', triple.key())
    return [[triple.votes_one_m, triple.votes_one_f, triple.votes_one_k],
            [triple.votes_two_m, triple.votes_two_f, triple.votes_two_k],
            [triple.votes_three_m, triple.votes_three_f, triple.votes_three_k]]


def _AddTripleVoteCounts(triple_key, votes):
  """Adds votes to a triple's vote count.

  This should be run in a transaction.

  Args:
    triple_key: (db.Key) the triple to update
    votes: ([str]) a 3-list of 'm', 'f', and 'k', corresponding to the votes
        for the 3 items in the triple, in order.
  """
  triple = models.Triple.get(triple_key)
  if triple.has_cached_votes:
    triple.votes_one_m += 1 if votes[0] == 'm' else 0
    triple.votes_one_f += 1 if votes[0] == 'f' else 0
    triple.votes_one_k += 1 if votes[0] == 'k' else 0
    triple.votes_two_m += 1 if votes[1] == 'm' else 0
    triple.votes_two_f += 1 if votes[1] == 'f' else 0
    triple.votes_two_k += 1 if votes[1] == 'k' else 0
    triple.votes_three_m += 1 if votes[2] == 'm' else 0
    triple.votes_three_f += 1 if votes[2] == 'f' else 0
    triple.votes_three_k += 1 if votes[2] == 'k' else 0
    triple.put()
  else:
    logging.warning('_AddTripleVoteCounts: Legacy Triple without vote counts:'
                    '%s', triple_key)


def _UpdateTripleVoteCounts(triple_key, new_counts):
  """Updates vote counts on the given triple.

  Args:
    triple: (Triple) triple to update
    new_counts: ([[int]]) These values are the new values for votes_one_m, ...,
        votes_three_k. See core.GetTripleVoteCounts.
  """
  triple = models.Triple.get(triple_key)

  assert(len(new_counts) == 3)
  votes_one, votes_two, votes_three = new_counts
  assert(len(votes_one) == 3)
  assert(len(votes_two) == 3)
  assert(len(votes_three) == 3)

  triple.votes_one_m, triple.votes_one_f, triple.votes_one_k = votes_one
  triple.votes_two_m, triple.votes_two_f, triple.votes_two_k = votes_two
  triple.votes_three_m, triple.votes_three_f, triple.votes_three_k = (
      votes_three)
  triple.has_cached_votes = True

  triple.put()


def MakeEntity(name, query, user_ip, thumb_url, original_url):
  """Makes an Entity with the given attributes."""
  # Get the thumbnail URL for the entity.  This could throw
  # URLError. We'll let it bubble up.
  result = urlfetch.fetch(thumb_url)
  logging.info('Downloading %s' % thumb_url)
  entity = models.Entity(name=name,
                         data=result.content,
                         query=query,
                         original_url=original_url)
  entity.put()

  return entity


def MakeTriple(entities, creator, creator_ip):
  """Create the named triple.

  Args:
    entities: a data structure built in MakeSubmitHandler.
    creator: the user who created the Triple.
    creator_ip: IP address of the request to make this triple.
  """
  for i in range(len(entities)):
    # TODO(mjkelly): Check for a signature element.
    for k in ['n', 'u', 'q', 'ou']:
      if not entities[i][k]:
        raise ValueError("Entity %s missing attribute '%s'" % (i, k))

  # This may raise a URLError or EntityValidatationError.
  one = MakeEntity(name=entities[0]['n'],
                   query=entities[0]['q'],
                   user_ip=creator_ip,
                   thumb_url=entities[0]['u'],
                   original_url=entities[0]['ou'])
  two = MakeEntity(name=entities[1]['n'],
                   query=entities[1]['q'],
                   user_ip=creator_ip,
                   thumb_url=entities[1]['u'],
                   original_url=entities[1]['ou'])
  three = MakeEntity(name=entities[2]['n'],
                     query=entities[2]['q'],
                     user_ip=creator_ip,
                     thumb_url=entities[2]['u'],
                     original_url=entities[2]['ou'])

  triple = models.Triple(one=one, two=two, three=three,
                         creator=creator,
                         creatorip=creator_ip,
                         has_cached_votes=True,
                         votes_one_m=0,
                         votes_one_f=0,
                         votes_one_k=0,
                         votes_two_m=0,
                         votes_two_f=0,
                         votes_two_k=0,
                         votes_three_m=0,
                         votes_three_f=0,
                         votes_three_k=0)
  triple.enable()
  triple.put()
  return triple


def MakeAssignment(triple_id, v1, v2, v3, user, user_ip):
  """Create a new assignment.

  Args:
    request: the POST request from the client
    user: the user who made the assignment request
  """
  values = [v1, v2, v3]
  if set(values) != set(['m', 'f', 'k']):
    return None

  try:
    triple_id = long(triple_id)
  except ValueError:
    logging.error("make_assignment: bad triple key '%s'", triple_id)

  triple = models.Triple.get_by_id(triple_id)
  logging.debug('triple = %s', triple)
  if triple is None:
    logging.error('make_assignment: No triple with key %s', triple_id)
    return None

  db.run_in_transaction(_AddTripleVoteCounts, triple.key(), values)

  # We get an entity->action map from the client, but we need to reverse
  # it to action->entity to update the DB.
  triple_entities = [triple.one,
                     triple.two,
                     triple.three]
  entities = {}
  for i in range(len(values)):
    # Items in values are guaranteed to be 'm', 'f', 'k' (check above)
    entities[values[i]] = triple_entities[i]

  if (entities['m'] is None or entities['f'] is None
      or entities['k'] is None):
    logging.error('Not all non-None: marry = %s, fuck = %s, kill = %s',
            entities['m'], entities['f'], entities['k'])
    return None

  assign = models.Assignment(triple=triple,
                             marry=entities['m'],
                             fuck=entities['f'],
                             kill=entities['k'],
                             user=user,
                             userip=str(user_ip))
  assign.put()
  logging.info("Assigned m=%s, f=%s, k=%s to %s", entities['m'],
      entities['f'], entities['k'], triple)
  return assign


def ImageSearch(query, user_ip):
  """Performs an image search. It should return 10 results.

  Args:
    query: (str) The search query
    user_ip: (str) IP address of user making the query, for accounting.

  Returns:
    [Image]: A list of Image objects representing search results
  """
  images = []
  query = query.encode('utf-8')
  url = ('https://www.googleapis.com/customsearch/v1'
         '?key={key}'
         '&cx={cx}'
         '&q={q}'
         '&userIp={userip}'
         '&searchType=image').format(
         key=config.CSE_API_KEY,
         cx=config.CSE_ID,
         q=urllib2.quote(query),
         userip=user_ip)
  logging.info('ImageSearch: query url=%s', url)

  download_start = datetime.datetime.now()
  # This may raise a DownloadError
  result = urlfetch.fetch(url)
  download_finish = datetime.datetime.now()
  data = json.loads(result.content)
  parse_finish = datetime.datetime.now()

  logging.info('ImageSearch: downloaded %s bytes; %s to download, %s to parse',
    len(result.content),
    download_finish - download_start,
    parse_finish - download_finish)

  for item in data['items']:
    link = item['link']
    thumb = item['image']['thumbnailLink']
    images.append(Image(original=link, thumbnail=thumb))

  return images

def Sign(*items):
  """Signs a sequence of items using our internal HMAC key.

  Args:
    *items: Any sequence of items.

  Returns:
    (str) Hex digest of *items
  """
  h = hmac.new(config.HMAC_KEY)
  for item in items:
    h.update(item)
  return h.hexdigest()

