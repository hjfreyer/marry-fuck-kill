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

from google.appengine.api import users
from google.appengine.ext import webapp
from google.appengine.ext import db

import models
import triple_handlers

def GetUserData(url_base):
  nickname = ''
  user = users.get_current_user()
  if user:
    nickname = user.nickname()

  return dict(nickname=nickname,
              login_url=users.create_login_url(url_base),
              logout_url=users.create_logout_url(url_base))


class MainPageHandler(webapp.RequestHandler):
  def get(self):
    rand = models.Triple.get_next_id(self.request, self.response)
    if not rand:
      self.redirect('/about')
    else:
      self.redirect('/vote/' + str(rand))


class AboutHandler(webapp.RequestHandler):
  def get(self):
    template_values = dict(page='about')
    template_values.update(GetUserData('/about'))
    template = ezt.Template('templates/about.html')
    template.generate(self.response.out, template_values)


class VoteHandler(webapp.RequestHandler):
  def get(self, triple_id):
    if not triple_id.isdigit():
      self.error(404)
      return

    triple = models.Triple.get_by_id(long(triple_id))

    if not triple:
      self.error(404)
      return

    prev_id = self.request.get('prev')
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

    template_values = dict(page='vote',
                           triple_id=triple_id,
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
    template = ezt.Template('templates/vote.html')
    template.generate(self.response.out, template_values)


class VoteSubmitHandler(webapp.RequestHandler):
  def post(self):
    action = self.request.get('action')

    logging.info('Vote handler. Action: %s', action)

    if action == 'submit':
      models.Assignment.make_assignment(self.request)
      query_suffix = '?prev=%s' % self.request.get('key')
    else:
      query_suffix = ''

    rand = models.Triple.get_next_id(self.request, self.response)
    self.redirect('/vote/%s%s' % (str(rand), query_suffix))


class MakeSubmitHandler(webapp.RequestHandler):
  def post(self):
    logging.info('Make handler')

    triple = triple_handlers.TripleCreationHandler.MakeTriple(
      self.request, users.get_current_user())

    self.redirect('/vote/%s?new' % triple.key().id())


class MakeHandler(webapp.RequestHandler):
  def get(self):
    template_values = dict(page='make')
    template_values.update(GetUserData('/make'))
    template = ezt.Template('templates/make.html')
    template.generate(self.response.out, template_values)


class MyMfksHandler(webapp.RequestHandler):
  def get(self):
    template_values = dict(page='mine')
    user = users.get_current_user()
    triples = [t for t in models.Triple.all().filter('creator =', user)]

    items = []
    for t in triples:
      item = EztItem(key=str(t.key().id()),
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
    template = ezt.Template('templates/mymfks.html')
    template.generate(self.response.out, template_values)

class GenerateRandHandler(webapp.RequestHandler):
  """Regenerates the 'rand' attributes of all Triples.

  We go through some acrobatics to batch our query to avoid running afoul of
  the 1000-item limit.
  """
  def get(self):
    batch_size = 1000
    template_values = dict(page='mine')
    user = users.get_current_user()
    triples = db.Query(models.Triple).order('__key__').fetch(batch_size)
    logging.info('GenerateRandHandler: initial query got %d', len(triples))
    count = 0
    while triples:
      for t in triples:
        count += 1
        t.rand = random.random()
        logging.debug('GenerateRandHandler: id=%s rand=%.15f', t.key(), t.rand)
        t.put()
      triples = db.Query(models.Triple).filter('__key__ >', t.key()).order('__key__').fetch(batch_size)
      logging.info('GenerateRandHandler: subsequent query got %d', len(triples))
    self.response.out.write('Regenerated %d Triples' % count)

class EztItem(object):
  """A simple wrapper to convert dict keys to object attributes.

  This makes sending complex objects to EZT much easier.
  """
  def __init__(self, **kwargs):
    vars(self).update(kwargs)
