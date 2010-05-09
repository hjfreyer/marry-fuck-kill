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

from google.appengine.ext import db

class Entity(db.PolyModel):
    name = db.StringProperty(required=True)
    creator = db.UserProperty()    

    type = db.StringProperty(required=True, choices=set(["abstract", "person", "object"]))

    def __str__(self):
        return "[Entity: %s]" % self.name

class Triple(db.PolyModel):
    name = db.StringProperty(required=True)

    creator = db.UserProperty()
 
    quality = db.FloatProperty()

    one = db.ReferenceProperty(Entity)
    two = db.ReferenceProperty(Entity)
    three = db.ReferenceProperty(Entity)

    def __str__(self):
        return "[Triple: %s, %s, %s]" % (self.one, self.two, self.three)
   
class Assignment(db.PolyModel):
    user = db.UserProperty()

    triple = db.ReferenceProperty(Triple)

    marry = db.ReferenceProperty(Entity)
    fuck = db.ReferenceProperty(Entity)
    kill = db.ReferenceProperty(Entity)
