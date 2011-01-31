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

import models

class MainPageHandler(webapp.RequestHandler):
  def get(self):
    rand = models.Triple.get_random_id()
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
    if not triple_id.isdigit():
      self.error(404)
      return

    triple = models.Triple.get_by_id(long(triple_id))

    if not triple:
      self.error(404)
      return

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
                                              e3_url=three.get_full_url()))


class VoteSubmitHandler(webapp.RequestHandler):
  def post(self):
    action = self.request.get('action')

    logging.error(action)

    if action == 'submit':
      models.Assignment.make_assignment(self.request)

    rand = models.Triple.get_random_id()
    self.redirect('/vote/' + str(rand))


class MakeHandler(webapp.RequestHandler):
  def get(self):
    template = ezt.Template('templates/make.html')
    template.generate(self.response.out, dict(page='make'))

