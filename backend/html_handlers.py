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

import datetime
import ezt
import json
import logging
import urllib2

from google.appengine.api import users
from google.appengine.ext import db
from google.appengine.ext import webapp

import core
import ezt_util
import models

def _LogRequest(page_name, req):
    """Logs a reqest object in a a relatively human-readable way.

    Args:
        page_name: A string identifying this page. Used in log entry for easy
            greppability.
        req: A WSGI request object.
    """
    arg_dict = dict([(arg, req.get(arg)) for arg in req.arguments()])
    logging.info('%s page params: %s', page_name, arg_dict)

def GetUserData(url_base):
  # This must be in this module to have access to the current user.
  nickname = ''
  user = users.get_current_user()
  if user:
    nickname = user.nickname()

  if users.is_current_user_admin():
    is_current_user_admin = True
  else:
    is_current_user_admin = None

  return dict(nickname=nickname,
              login_url=users.create_login_url(url_base),
              logout_url=users.create_logout_url(url_base),
              is_current_user_admin=is_current_user_admin)


class RequestHandler(webapp.RequestHandler):
  def error(self, code):
    super(RequestHandler, self).error(code)
    if code == 404:
      template_values = dict(page='')
      template_values.update(GetUserData('/about'))
      ezt_util.WriteTemplate('notfound.html', template_values,
                             self.response.out)


class MainPageHandler(RequestHandler):
  def get(self):
    rand = models.Triple.get_next_id(self.request, self.response)
    if not rand:
      self.redirect('/about')
    else:
      RenderVotePage(self, str(rand))


class AboutHandler(RequestHandler):
  def get(self):
    template_values = dict(page='about')
    template_values.update(GetUserData('/about'))
    ezt_util.WriteTemplate('about.html', template_values, self.response.out)


class VoteHandler(RequestHandler):
  def get(self, triple_id):
    RenderVotePage(self, triple_id)


def RenderVotePage(handler, triple_id):
  if not triple_id.isdigit():
    handler.error(404)
    return

  triple = models.Triple.get_by_id(long(triple_id))

  if not triple:
    handler.error(404)
    return

  prev_id = handler.request.get('prev')
  display_prev_id = None
  prev_names = ['', '', '']
  prev_urls = ['', '', '']
  logging.info('Vote page for %s. Prev = %s', triple_id, prev_id)

  if prev_id:
    prev_triple = models.Triple.get_by_id(int(prev_id))
    if prev_triple is not None:
      prev_entities = [prev_triple.one, prev_triple.two, prev_triple.three]
      prev_names = [e.name for e in prev_entities]
      prev_urls = core.GetStatsUrlsForTriple(prev_triple)
      display_prev_id = prev_id

  one = triple.one
  two = triple.two
  three = triple.three

  template_values = dict(page='vote',
                         triple=triple,
                         triple_id=triple_id,
                         e1_name=one.name.encode('utf-8'),
                         e1_url=one.image_url,
                         e2_name=two.name.encode('utf-8'),
                         e2_url=two.image_url,
                         e3_name=three.name.encode('utf-8'),
                         e3_url=three.image_url,
                         prev_id=display_prev_id,
                         prev_e1_name=prev_names[0].encode('utf-8'),
                         prev_e2_name=prev_names[1].encode('utf-8'),
                         prev_e3_name=prev_names[2].encode('utf-8'),
                         prev_e1_stat_url=prev_urls[0],
                         prev_e2_stat_url=prev_urls[1],
                         prev_e3_stat_url=prev_urls[2])
  template_values.update(GetUserData('/vote/' + triple_id))
  ezt_util.WriteTemplate('vote.html', template_values, handler.response.out)


class VoteSubmitHandler(RequestHandler):
  def post(self):
    _LogRequest('Vote', self.request)
    action = self.request.get('action')

    if action == 'submit':
      core.MakeAssignment(triple_id=self.request.get('key'),
                          v1=self.request.get('v1'),
                          v2=self.request.get('v2'),
                          v3=self.request.get('v3'),
                          user=users.get_current_user(),
                          user_ip=self.request.remote_addr)

    query_suffix = '?prev=%s' % self.request.get('key')

    rand = models.Triple.get_next_id(self.request, self.response)
    self.redirect('/vote/%s%s' % (str(rand), query_suffix))


class MakeSubmitHandler(RequestHandler):
  def post(self):
    """Handles a request to make a new Triple.

    We expect the following parameters:
      n1, n2, n3: Names of the new triples
      u1, u2, u3: The URLs of the thumbnails of the new triples.
      q1, q2, q3: The queries that generated the new trples.
      ou1, ou2, ou3: The original URLs of the new triples.
    """
    # TODO(mjkelly): When we have a new client, check the 'sig' values we get.
    # That will allow us to avoid repeating the search on the server side.
    _LogRequest('Make', self.request)

    entities = []
    for n in range(1, 4):
      entities.append({'n': self.request.get('n' + str(n)),
                       'u': self.request.get('u' + str(n)),
                       'q': self.request.get('q' + str(n)),
                       'ou': self.request.get('ou' + str(n))})

    triple = core.MakeTriple(entities,
                             creator=users.get_current_user(),
                             creator_ip=self.request.remote_addr)
    self.redirect('/vote/%s?new' % triple.key().id())


class MakeHandler(RequestHandler):
  def get(self):
    template_values = dict(page='make')
    template_values.update(GetUserData('/make'))
    ezt_util.WriteTemplate('maker.html', template_values, self.response.out)


class EntityImageHandler(RequestHandler):
  def get(self, entity_id):
    try:
      entity = models.Entity.get(urllib2.unquote(entity_id))
    except db.BadKeyError:
      self.error(404)
      return
    self.response.headers['Content-Type'] = 'image/jpg';
    self.response.out.write(entity.data)


class MyMfksHandler(RequestHandler):
  def get(self):
    template_values = dict(page='mymfks')
    user = users.get_current_user()
    query = models.Triple.all().filter('creator =', user).order('time')
    triples = [t for t in query]

    items = []
    for t in triples:
      stats = core.GetStatsUrlsForTriple(t)
      item = ezt_util.EztItem(key=str(t.key().id()),
                              triple=t,
                              one_stats=stats[0],
                              two_stats=stats[1],
                              three_stats=stats[2])
      items.append(item)

    template_values.update(GetUserData('/mymfks'))
    template_values.update(dict(triples=items))
    ezt_util.WriteTemplate('mymfks.html', template_values, self.response.out)


class ImageSearchHandler(RequestHandler):
  def get(self):
    _LogRequest('ImageSearch', self.request)
    query = self.request.get('q')

    images = core.ImageSearch(query, self.request.remote_addr)
    images_dicts = []

    for img in images:
      d = img._asdict()
      # The goal here is just to ensure that we can check that we once returned
      # this URL as a result for a search. We're not attempting to associate
      # the URL with specific search terms, or with a time period.
      #
      # If someone puts in enough effort to abuse this somehow, we can think
      # about locking it down more.
      d.update(dict(sig=core.Sign(*img)))
      images_dicts.append(d)

    results = json.dumps(dict(time=str(datetime.datetime.now()), images=images_dicts))

    logging.info('results = %s', images)
    self.response.headers['Content-Type'] = 'application/json';
    self.response.out.write(results)


class CatchAllHandler(RequestHandler):
  def get(self):
    self.error(404)
