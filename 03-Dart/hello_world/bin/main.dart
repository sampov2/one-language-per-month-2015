// Copyright (c) 2015, Sampo Savolainen. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

import 'package:hello_world/hello_world.dart' as hello_world;
import 'package:hello_world/yikes.dart' as yikes;

main() {
  // Libraries are singleton, as proved by how the value of yikes.testit() in increased
  // by the call to hello_world.calculate()
  print('yikes.testit(): ${yikes.testit()}!');
  print('Hello world: ${hello_world.calculate()}!');
  print('yikes.testit(): ${yikes.testit()}!');
}
