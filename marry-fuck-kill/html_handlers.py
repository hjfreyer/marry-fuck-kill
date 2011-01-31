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

from google.appengine.ext import webapp

import assignment_handlers
import models

class MainPageHandler(webapp.RequestHandler):
  def get(self):
    rand = models.Triple.get_random_key()
    if not rand:
      self.redirect('/about')
    else:
      self.redirect('/vote/' + str(rand))


class AboutHandler(webapp.RequestHandler):
  def get(self):
    template = ezt.Template('templates/about.html')
    template.generate(self.response.out, dict(page='about'))


class VoteHandler(webapp.RequestHandler):
  def get(self, triple_id):
    triple = models.Triple.get(triple_id)

    prev_id = self.request.get('prev')
    logging.info('Vote page for %s. Prev = %s', triple_id, prev_id)

    if prev_id:
      prev_triple = models.Triple.get(prev_id)
      prev_entities = [prev_triple.one, prev_triple.two, prev_triple.three]
      prev_names = [e.name for e in prev_entities]
      prev_urls = [e.get_stats_url() for e in prev_entities]
    else:
      prev_names = ['', '', '']
      prev_urls = ['', '', '']

    one = triple.one
    two = triple.two
    three = triple.three

    template = ezt.Template('templates/vote.html')
    template.generate(self.response.out, dict(page='vote',
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
                                              prev_e3_stat_url=prev_urls[2]))


class VoteSubmitHandler(webapp.RequestHandler):
  def post(self):
    action = self.request.get('action')

    logging.info('Vote handler. Action: %s', action)

    if action == 'submit':
      assignment_handlers.AssignmentHandler.make_assignment(self.request)
      query_string = '?prev=%s' % self.request.get('key')
    else:
      query_string = ''

    rand = models.Triple.get_random_key()
    self.redirect('/vote/%s%s' % (str(rand), query_string))


class MakeHandler(webapp.RequestHandler):
  def get(self):
    template = ezt.Template('templates/make.html')
    template.generate(self.response.out, dict(page='make'))

