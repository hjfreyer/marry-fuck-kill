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

from google.appengine.ext import db

class Entity(db.Model):
    name = db.StringProperty(required=True)
    creator = db.UserProperty()    

    type = db.StringProperty(choices=set(["abstract", "person", "object"]))
   
    def __init__(self, *args, **kwargs):
        key = Entity.get_key(**kwargs)
        kwargs["key_name"] = key
        super(Entity, self).__init__(*args, **kwargs)
        self.key_name =  key

    def __str__(self):
        return "[Entity: %s]" % self.name

    @staticmethod
    def get_key(name **kwargs):
        return "%s" % name
    
class Triple(db.Model):
    creator = db.UserProperty()
 
    quality = db.FloatProperty(default=1.0)

    one = db.ReferenceProperty(Entity,
                               collection_name="triple_reference_one_set")
    two = db.ReferenceProperty(Entity,
                               collection_name="triple_reference_two_set")
    three = db.ReferenceProperty(Entity,
                                 collection_name="triple_reference_three_set")

    def __init__(self, *args, **kwargs):
        key = Triple.get_key(**kwargs)
        kwargs["key_name"] = key
        super(Triple, self).__init__(*args, **kwargs)
        self.key_name = key

    def __str__(self):
        return "[Triple: %s, %s, %s]" % (self.one, self.two, self.three)

    @staticmethod
    def get_key(one, two, three, **kwargs):
        return "%s.%s.%s" % (one.key_name, two.key_name, three.key_name)
    
class Assignment(db.Model):
    user = db.UserProperty()

    triple = db.ReferenceProperty(Triple)

    marry = db.ReferenceProperty(Entity,
                                 collection_name="assignment_reference_marry_set")
    fuck = db.ReferenceProperty(Entity,
                                 collection_name="assignment_reference_fuck_set")
    kill = db.ReferenceProperty(Entity,
                                 collection_name="assignment_reference_kill_set")
