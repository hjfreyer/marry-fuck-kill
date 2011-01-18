#!/usr/bin/env python
#
# Copyright 2010 Hunter Freyer and Michael Kelly
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

import logging
import random

from google.appengine.ext import db

class EntityValidationError(Exception):
    """There was an error validating an entity."""
    pass

class Entity(db.Model):
    name = db.StringProperty(required=True)
    url = db.StringProperty(required=False)
    query = db.StringProperty(required=False)
    creator = db.UserProperty()    

    type = db.StringProperty(choices=set(["abstract", "person", "object"]))
   
    def __str__(self):
        return self.name

    def __repr__(self):
        return "Entity(%r)" % self.name

    def json(self):
        return {'name': self.name, 'url': self.get_full_url()}

    def get_full_url(self):
        """Get the full URL including the prefix."""
        return self.url

    def set_full_url(self, url):
        """Given a full URL, set the url property."""
        logging.info('set_full_url: url=%s', url)
        self.url = url
        

def PutEntity(name, url, query):
    entity = Entity(name=name, query=query)
    entity.set_full_url(url)
    entity.put()
    return entity
    
class Triple(db.Model):
    creator = db.UserProperty()
 
    quality = db.FloatProperty(default=1.0)

    one = db.ReferenceProperty(Entity,
                               collection_name="triple_reference_one_set")
    two = db.ReferenceProperty(Entity,
                               collection_name="triple_reference_two_set")
    three = db.ReferenceProperty(Entity,
                                 collection_name="triple_reference_three_set")
        
    @staticmethod
    def get_random():
        keys = [k for k in db.Query(Triple, keys_only=True).filter('quality >', 0.0)]
        if not keys:
            return None
        else:
            return Triple.get(random.choice(keys))

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
    def validate(one, two, three):
        """Verify that the 3 given entities can make a valid triple.

        Raises EntityValidationError if there was a problem.
        """

        if len(set([one.name, two.name, three.name])) < 3:
            raise EntityValidationError("All item names must be distinct.")
        if len(set([one.url, two.url, three.url])) < 3:
            raise EntityValidationError("All item URLs must be distinct.")
        # TODO(mjkelly): implement me
        #raise EntityValidationError("test error")
        pass

def PutTriple(one, two, three):
    """Put a triple in the DB with canonical key.

    See Triple.key_name_from_entities().
    """
    triple = Triple(one=one, 
                    two=two,
                    three=three)
    triple.put()
    return triple

class Assignment(db.Model):
    user = db.UserProperty()

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
