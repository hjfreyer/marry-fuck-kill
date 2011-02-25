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

import ezt
import logging
import random
import urllib2

from google.appengine.api import users
from google.appengine.ext import db
from google.appengine.ext import webapp

import ezt_util
import models


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
  logging.info('Vote page for %s. Prev = %s', triple_id, prev_id)

  if prev_id:
    prev_triple = models.Triple.get_by_id(int(prev_id))
    prev_entities = [prev_triple.one, prev_triple.two, prev_triple.three]
    prev_names = [e.name for e in prev_entities]
    prev_urls = [e.get_stats_url() for e in prev_entities]
  else:
    prev_names = ['', '', '']
    prev_urls = ['', '', '']

  one = triple.one
  two = triple.two
  three = triple.three

  # Map False -> None so EZT understands.
  if triple.is_enabled():
    triple_enabled = True
  else:
    triple_enabled = None

  template_values = dict(page='vote',
                         triple_id=triple_id,
                         triple_rand=triple.rand,
                         triple_enabled=triple_enabled,
                         triple_creator=str(triple.creator),
                         triple_creatorip=str(triple.creatorip),
                         triple_reviewed=triple.reviewed,
                         triple_time=str(triple.time),
                         e1_name=one.name,
                         e1_url=one.get_full_url(),
                         e2_name=two.name,
                         e2_url=two.get_full_url(),
                         e3_name=three.name,
                         e3_url=three.get_full_url(),
                         prev_id=prev_id,
                         prev_e1_name=prev_names[0],
                         prev_e2_name=prev_names[1],
                         prev_e3_name=prev_names[2],
                         prev_e1_stat_url=prev_urls[0],
                         prev_e2_stat_url=prev_urls[1],
                         prev_e3_stat_url=prev_urls[2])
  template_values.update(GetUserData('/vote/' + triple_id))
  ezt_util.WriteTemplate('vote.html', template_values, handler.response.out)


class VoteSubmitHandler(RequestHandler):
  def post(self):
    action = self.request.get('action')

    logging.info('Vote handler. Action: %s', action)

    if action == 'submit':
      models.Assignment.make_assignment(self.request, users.get_current_user())

    query_suffix = '?prev=%s' % self.request.get('key')

    rand = models.Triple.get_next_id(self.request, self.response)
    self.redirect('/vote/%s%s' % (str(rand), query_suffix))


class MakeSubmitHandler(RequestHandler):
  def post(self):
    logging.info('Make handler')
    triple = models.Triple.make_triple(self.request, users.get_current_user())
    self.redirect('/vote/%s?new' % triple.key().id())


class MakeHandler(RequestHandler):
  def get(self):
    template_values = dict(page='make')
    template_values.update(GetUserData('/make'))
    ezt_util.WriteTemplate('make.html', template_values, self.response.out)


class EntityImageHandler(RequestHandler):
  def get(self, entity_id):
    entity = models.Entity.get(urllib2.unquote(entity_id))
    self.response.headers['Content-Type'] = "image/jpg";
    self.response.out.write(entity.data)


class MyMfksHandler(RequestHandler):
  def get(self):
    template_values = dict(page='mymfks')
    user = users.get_current_user()
    triples = [t for t in models.Triple.all().filter('creator =', user)]

    items = []
    for t in triples:
      item = ezt_util.EztItem(key=str(t.key().id()),
                              one_name=t.one.name,
                              two_name=t.two.name,
                              three_name=t.three.name,
                              one_url=t.one.get_full_url(),
                              two_url=t.two.get_full_url(),
                              three_url=t.three.get_full_url(),
                              one_stats=t.one.get_stats_url(),
                              two_stats=t.two.get_stats_url(),
                              three_stats=t.three.get_stats_url())
      items.append(item)

    template_values.update(GetUserData('/mymfks'))
    template_values.update(dict(triples=items))
    ezt_util.WriteTemplate('mymfks.html', template_values, self.response.out)


class CatchAllHandler(RequestHandler):
  def get(self):
    self.error(404)
