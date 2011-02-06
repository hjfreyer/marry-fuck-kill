#!/usr/bin/env python
#
# Copyright 2010 Hunter Freyer and Michael Kelly
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

import logging
import random
import urllib2
from django.utils import simplejson

from google.appengine.api import users
from google.appengine.ext import db

class EntityValidationError(Exception):
  """There was an error validating an entity."""
  pass

class Entity(db.Model):
  name = db.StringProperty(required=True)
  data = db.BlobProperty(required=True)
  query = db.StringProperty(required=False)
  creator = db.UserProperty()

  # Where images live on the server.
  BASE_URL = '/i/'

  # TODO(mjkelly): do something with this or get rid of it
  type = db.StringProperty(choices=set(['abstract', 'person', 'object']))

  def __str__(self):
    return self.name

  def __repr__(self):
    return 'Entity(%r)' % self.name

  def json(self):
    return {'name': self.name, 'url': self.get_full_url()}

  def get_full_url(self):
    logging.debug('get_full_url (for %d): returning %s',
        self.key().id(), Entity.BASE_URL + str(self.key()))
    return Entity.BASE_URL + str(self.key())

  def get_stats_url(self, w=160, h=85):
    """Returns the URL to a chart of MFK stats for this Entity.

    Args:
      w: (int) image width
      h: (int) image height
    """
    m = self.assignment_reference_marry_set.count()
    f = self.assignment_reference_fuck_set.count()
    k = self.assignment_reference_kill_set.count()

    url = ('http://chart.apis.google.com/chart'
        '?chxr=0,0,%(max)d'
        '&chxt=y'
        '&chbh=a'
        '&chs=%(w)dx%(h)d'
        '&cht=bvg'
        '&chco=9911BB,C76FDD,63067A'
        '&chds=0,%(max)d,0,%(max)d,0,%(max)d'
        '&chd=t:%(m)d|%(f)d|%(k)d'
        '&chdl=Marry|Fuck|Kill'
        '&chdlp=r' % (dict(m=m, f=f, k=k, max=max([m,f,k]), w=w, h=h)))
    return url

  @staticmethod
  def make_entity(name, url, query):
    if not url.startswith('http://images.google.com/images?q'):
      raise EntityValidationError('URL must come from Google image search.')

    # This could throw URLError. We'll let it bubble up.
    fh = urllib2.urlopen(url)
    logging.info('Downloading %s' % url)
    data = fh.read()
    logging.info('Downloaded %s bytes' % len(data))

    entity = Entity(name=name, data=data, query=query)

    entity.put()
    return entity


class Triple(db.Model):
  # User who created the Triple, or None if no user was logged in.
  creator = db.UserProperty()
  # These fields are not required for DB backwards-compatibility.
  creatorip = db.StringProperty(required=False, default='')
  time = db.DateTimeProperty(required=False, auto_now_add=True)
  # A measure of quality. Currently ignored.
  # TODO(mjkelly): Remove this.
  quality = db.FloatProperty(default=1.0)
  # A random value used to determine a randomized, repeatable display order.
  # Not necessarily unique. 'default' value is for backwards-compatibility.
  # If this is None, this Triple will not be returned by get_next_id() under
  # any cirumstances.
  rand = db.FloatProperty(required=False, default=0.0)

  one = db.ReferenceProperty(Entity,
                 collection_name="triple_reference_one_set")
  two = db.ReferenceProperty(Entity,
                 collection_name="triple_reference_two_set")
  three = db.ReferenceProperty(Entity,
                 collection_name="triple_reference_three_set")

  CONTEXT_COOKIE_NAME = 'mfkcontext'

  GOOGLE_API_KEY = 'ABQIAAAA4AIACTDq7g0UgEEe0e4XcBScM50iuTtmL4hn6SVBcuHEk5GnyBRYi46EgwfJeghlh-_jWgC9BbPapQ'
  SEARCH_REFERER = 'http://marry-fuck-kill.appspot.com'

  def disable(self):
    """Prevent this Triple from being picked in the random rotation.

    It will still be visible from a direct link.

    This works by removing the 'rand' attribute from this object entirely, to
    prevent it from being returned in sort operations involving 'rand'.
    """
    self.rand = -1.0

  def enable(self):
    """Allow this Triple to be picked in the random rotation."""
    self.rand = random.random()

  @staticmethod
  def get_next_id(request, response):
    """Gets the next Triple ID, apparently at random.

    We use cookies in the request, if available, to determine what the next
    Triple will be.

    Args:
      request: request object from the client
      response: the response object that will be sent to the client.

    Returns:
      (ID) The ID of the next Triple, and a cookie to set on the client
           to maintain state.
    """
    # We _could_ also store the random value on the client side if we want to
    # cut down on the number of DB requests we make.
    prev_id = None
    next_id = None
    try:
      # This fails if there is no cookie, or it isn't a valid long.
      prev_id = long(request.cookies.get(Triple.CONTEXT_COOKIE_NAME))
    except TypeError:
      pass
    logging.info('get_next_id: prev_id = %s', prev_id)
    if prev_id is not None:
      prev_triple = Triple.get_by_id(long(prev_id))
      if prev_triple is not None:
        prev_rand = prev_triple.rand
        next_id = Triple._get_greater_rand(prev_rand)

    # If we failed to intelligently determine the next triple, pick at random.
    if next_id is None:
      next_id = Triple._get_random_id()

    # TODO(mjkelly): Try to find a lightweight wrapper for cookie-setting
    # instead of manually constructing them.
    if next_id is not None:
      response.headers.add_header('Set-Cookie', '%s=%d' % (Triple.CONTEXT_COOKIE_NAME, next_id))

    return next_id

  @staticmethod
  def _get_greater_rand(rand, _or_equal=False):
    """Get the Triple ID with the next-greater rand value.

    This method contains all the real logic for picking random values.

    Args:
      rand: The value to compare against.
    """
    if _or_equal:
      logging.info('_get_greater_rand: rand >= %s', rand)
      query = db.Query(Triple, keys_only=True).filter(
          'rand >=', rand).order('rand').order('__key__')
    else:
      logging.info('_get_greater_rand: rand > %s', rand)
      query = db.Query(Triple, keys_only=True).filter(
          'rand >', rand).order('rand').order('__key__')
    if not query.count():
      if not _or_equal:
        # If we were using a rand value > 0, we probably hit the end of the
        # list naturally. We want to wrap around to the lowest triple, so we
        # use rand >= 0.0.
        return Triple._get_greater_rand(0.0, _or_equal=True)
      else:
        # If we returned no results and yet were using a rand value < 0, the
        # database is empty or contains only invalid rand values. Give up. This
        # is the base case for the recursion above.
        return None
    else:
      return query.get().id()

  @staticmethod
  def _get_random_id():
    """Returns a truly random Triple.

    This is a fallback used by get_next_id for when we don't have any context
    on previously-seen Triples from the client.
    """
    rand = random.random()
    logging.info('_get_random_id: rand = %f', rand)
    return Triple._get_greater_rand(rand)

  def __str__(self):
    return Triple.name_from_entities(self.one, self.two, self.three)

  def __repr__(self):
    return "Triple(one=%r, two=%r, three=%r)" % (self.one, self.two, self.three)

  def json(self):
    return {'one': self.one.json(),
        'two': self.two.json(),
        'three': self.three.json(),
        'key': str(self.key())}

  @staticmethod
  def name_from_entities(one, two, three):
    """Given three Entities, generate the canonical key name for a Triple
    containing them.
    """
    keys = [one.name, two.name, three.name]
    keys.sort()
    return "%s.%s.%s" % (keys[0], keys[1], keys[2])

  @staticmethod
  def make_triple(request, creator):
    """Create the named triple.

    Args:
      request: the POST request from the client
      creator: the user who created the Triple.

    The only non-obvious part of this is that we check that q[1-3] actually
    include u[1-3]. This is to prevent users from adding any URL they
    please.
    """
    entities = Triple.parse_request(request)
    for i in range(len(entities)):
      for k in ['n', 'u', 'q']:
        if not entities[i][k]:
          raise ValueError("Entity %s missing attribute '%s'" % (i, k))

    # This may raise an EntityValidationError.
    Triple.validate_request(entities, request.remote_addr)

    # This may raise a URLError or EntityValidatationError.
    one = Entity.make_entity(entities[0]['n'], entities[0]['u'],
                             entities[0]['q'])
    two = Entity.make_entity(entities[1]['n'], entities[1]['u'],
                             entities[1]['q'])
    three = Entity.make_entity(entities[2]['n'], entities[2]['u'],
                               entities[2]['q'])

    triple = Triple(one=one, two=two, three=three,
                    creator=creator,
                    creatorip=request.remote_addr,
                    rand=random.random())
    triple.put()
    return triple

  @staticmethod
  def validate_request(req, userip):
    """Verify that the 3 given entities can make a valid triple.


    Args:
      req_dict: (list) The data structure from Triple.parse_request.
      userip: (str) The IP of the user who made the request

    Raises:
      EntityValidationError if the triple does not validate.
      DownloadError if there was a network error checking the triple. 
    """
    # This must be synchronized with the number of pages the creation interface
    # shows. Raising it incurs a large performance penalty.
    check_pages = 3

    names = [item['n'] for item in req]
    urls = [item['u'] for item in req]
    if len(set(names)) < 3:
      raise EntityValidationError('All item names must be distinct: %s', names)
    if len(set(urls)) < 3:
      raise EntityValidationError('All item URLs must be distinct: %s', urls)

    for item in req:
      # This may raise a DownloadError
      images = Triple._get_images_for_query(item['q'], check_pages, userip)
      search_urls = [image['tbUrl'] for image in images]
      if item['u'] not in search_urls:
        logging.error("URL '%s' is not in result set for query '%s'. "
                      "Result set over %d pages is: %s" % (
                          item['u'], item['q'], check_pages, search_urls))
        raise EntityValidationError(
            "URL '%s' is not in result set for query '%s'." % (item['u'],
                                                               item['q']))

  @staticmethod
  def _get_images_for_query(query, num_pages, userip):
    start = 0
    images = []
    for page in range(num_pages):
      url = ('https://ajax.googleapis.com/ajax/services/search/images'
             '?v=1.0'
             '&safe=moderate'
             '&userip=%(ip)s'
             '&q=%(q)s'
             '&start=%(start)s'
             '&key=%(key)s' % {'q': urllib2.quote(query),
                               'start': start,
                               'ip': userip,
                               'key': Triple.GOOGLE_API_KEY})
      logging.info('_get_images_for_query: url=%s', url)
      req = urllib2.Request(url, None, {'Referer': Triple.SEARCH_REFERER})
      # This may raise a DownloadError
      results = simplejson.load(urllib2.urlopen(req))
      images += results['responseData']['results']
      logging.info('_get_images_for_query: %d results so far', len(images))

      # Make sure there are more pages before advancing 'start'
      pages = results['responseData']['cursor']['pages']
      if len(pages) >= num_pages:
        start = pages[page+1]['start']
      else:
        break

    return images

  @staticmethod
  def parse_request(request):
    """Parses a Triple creation request and puts it in a convenient format.

    We expect the following request params:
    n[1-3]: the 3 triple display names
    u[1-3]: the 3 triple image URLs
    q[1-3]: the search string used to find u[1-3]
    """
    return [{'n': request.get('n1'),
             'u': request.get('u1'),
             'q': request.get('q1')},
            {'n': request.get('n2'),
             'u': request.get('u2'),
             'q': request.get('q2')},
            {'n': request.get('n3'),
             'u': request.get('u3'),
             'q': request.get('q3')}]


class Assignment(db.Model):
  user = db.UserProperty()
  # These fields are not required for DB backwards-compatibility.
  userip = db.StringProperty(required=False, default='')
  time = db.DateTimeProperty(required=False, auto_now_add=True)

  triple = db.ReferenceProperty(Triple)

  marry = db.ReferenceProperty(Entity,
                 collection_name="assignment_reference_marry_set")
  fuck = db.ReferenceProperty(Entity,
                 collection_name="assignment_reference_fuck_set")
  kill = db.ReferenceProperty(Entity,
                 collection_name="assignment_reference_kill_set")

  def __str__(self):
    return "Assignment(marry=%s, fuck=%s, kill=%s)" % (
        self.marry.key().name(),
        self.fuck.key().name(),
        self.kill.key().name())

  def __repr__(self):
    return str(self)

  @staticmethod
  def make_assignment(request, user):
    """Create a new assignment.

    Args:
      request: the POST request from the client
      user: the user who made the assignment request
    """
    triple_key = request.get('key')
    values = [request.get('v1'), request.get('v2'), request.get('v3')]

    logging.debug("triple_key = %s", triple_key)
    logging.debug("values = %s", values)

    if set(values) != set(['m', 'f', 'k']):
      return None

    # We get an entity->action map from the client, but we need to reverse
    # it to action->entity to update the DB.
    triple = Triple.get_by_id(int(triple_key))
    logging.debug("triple = %s", triple)
    if triple is None:
      logging.error("No triple with key %s", triple_key)
      return None

    triple_entities = [triple.one,
                       triple.two,
                       triple.three]
    entities = {}
    for i in range(len(values)):
      # Items in values are guaranteed to be 'm', 'f', 'k' (check above)
      entities[values[i]] = triple_entities[i]

    if (entities['m'] is None or entities['f'] is None
        or entities['k'] is None):
      logging.error("Not all non-None: marry = %s, fuck = %s, kill = %s",
              entities['m'], entities['f'], entities['k'])
      return None

    assign = Assignment(marry=entities['m'],
                        fuck=entities['f'],
                        kill=entities['k'],
                        user=user,
                        userip=str(request.remote_addr))
    assign.put()
    logging.info("Assigned m=%s, f=%s, k=%s to %s", entities['m'],
        entities['f'], entities['k'], triple)
    return assign
