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
    url = db.StringProperty(required=False)
    creator = db.UserProperty()    

    # TODO(mjk): Consider whether this is really necessary. Right now it's
    # cheap paranoia.
    URL_BASE = 'http://images.google.com/images?q='
    
    type = db.StringProperty(choices=set(["abstract", "person", "object"]))
   
    def __str__(self):
        return self.name

    def __repr__(self):
        return "Entity(%r)" % self.name

    def json(self):
        return {'name': self.name}

    def get_full_url(self):
        """Get the full URL including the prefix."""
        return self.URL_BASE + self.url

    def set_full_url(self, url):
        """Given a full URL, set the url property.

        We raise a ValueError if the URL does not begin with
        Entity.URL_BASE.
        """
        if url.startswith(self.URL_BASE):
            self.url = url[len(self.URL_BASE):]
        else:
            raise ValueError(url)
        

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
        keys = [k for k in db.Query(Triple, keys_only=True).filter('quality >', 0.0)]
        if not keys:
            return None
        else:
            return Triple.get(random.choice(keys))

    def __str__(self):
        return Triple.key_name_from_entities(self.one, self.two, self.three)

    def __repr__(self):
        return "Triple(one=%r, two=%r, three=%r)" % (self.one, self.two, self.three)

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
        return "Assignment(marry=%s, fuck=%s, kill=%s)" % (
                self.marry.key().name(),
                self.fuck.key().name(),
                self.kill.key().name())

    def __repr__(self):
        return str(self)
