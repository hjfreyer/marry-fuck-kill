#!/usr/bin/env python
#
# Copyright 2010 Hunter Freyer and Michael Kelly Inc.
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

import random

from google.appengine.ext import db

class Entity(db.Model):
    name = db.StringProperty(required=True)
    creator = db.UserProperty()    

    type = db.StringProperty(choices=set(["abstract", "person", "object"]))
   
    def __str__(self):
        return "[Entity: %s]" % self.name

    def json(self):
        return {'name': self.name}

def PutEntity(name):
    entity = Entity(name=name,
                    key_name=name)
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
        # TOOD(mjkelly): Make this constant, not O(triples).
        keys = []
        for k in Triple.all(keys_only=True):
            keys.append(k)

        if not keys:
            return None
        else:
            return Triple.get(random.choice(keys))

    def __str__(self):
        return "[Triple: %s, %s, %s]" % (self.one, self.two, self.three)

    def json(self):
        return {'one': self.one.json(),
                'two': self.two.json(),
                'three': self.three.json(),
                'key': self.key().name()}

    @staticmethod
    def key_name_from_entities(one, two, three):
        """Given three Entities, generate the canonical key name for a Triple
        containing them.
        """
        keys = [one.key().name(), two.key().name(), three.key().name()]
        keys.sort()
        return "%s.%s.%s" % (keys[0], keys[1], keys[2])

def PutTriple(one, two, three):
    """Put a triple in the DB with canonical key.

    See Triple.key_name_from_entities().
    """
    triple = Triple(one=one, 
                    two=two,
                    three=three,
                    key_name=Triple.key_name_from_entities(one, two, three))
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
        return "[Assignment: %s, %s, %s]" % (self.marry, self.fuck, self.kill)
