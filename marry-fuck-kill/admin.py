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

import ezt

from google.appengine.ext import webapp
from google.appengine.ext.webapp import util
from mapreduce import control

import ezt_util
import models

class MapReduceTriggerHandler(webapp.RequestHandler):
  MR_SPECS = dict(
    RANDOMIZE = dict(
      name='Randize',
      handler_spec='mapreduces.RandomizeMapper',
      reader_spec='mapreduce.input_readers.DatastoreInputReader',
      reader_parameters=dict(entity_kind='models.Triple'))
    )

  def get(self, mr_name):
    self.response.headers['Content-Type'] = 'text/plain'

    if mr_name in MapReduceTriggerHandler.MR_SPECS:
      control.start_map(**MapReduceTriggerHandler.MR_SPECS[mr_name])
      self.response.out.write('MR Cron Started: %s' % mr_name)
    else:
      self.response.out.write('MR Cron Not Found: %s' % mr_name)


class TripleReviewHandler(webapp.RequestHandler):
  """Admin-only handler to remove manually review Triples."""
  def get(self):
    ten = models.Triple.all().filter('reviewed =', False).fetch(10)

    ezt_util.WriteTemplate('review.html', dict(triples=ten), self.response.out)

  def post(self):
    ids = self.request.get('ids').split(',')
    ids = [long(i) for i in ids if i]

    triples = models.Triple.get_by_id(ids)

    for triple in triples:
      triple.reviewed = True
      triple.put()

    self.redirect('/admin/review')

class EnableDisableTripleHandler(webapp.RequestHandler):
  """Admin-only handler to remove Triple from random display."""
  def post(self):
    action = self.request.get('action')
    triple_id = self.request.get('key')
    triple = models.Triple.get_by_id(long(triple_id))
    if action == "disable":
      triple.disable()
    elif action == "enable":
      triple.enable()
    else:
      raise ValueError("Invalid action '%s'." % action)
    triple.put()

    self.redirect('/vote/%s' % triple_id)


def main():
  application = webapp.WSGIApplication([
      ("/admin/trigger_mr/(.*)", MapReduceTriggerHandler),
      ("/admin/enable_disable", EnableDisableTripleHandler),
      ("/admin/review", TripleReviewHandler),
    ], debug=True)
  util.run_wsgi_app(application)

if __name__ == '__main__':
  main()
