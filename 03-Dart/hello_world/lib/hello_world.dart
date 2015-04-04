// Copyright (c) 2015, <your name>. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

/// The hello_world library.
library hello_world;

int n = 10;

int calculate() {
  n++;
  return 6 * 7 * n;
}


int getN() {
  return n;
}