library tests_col;

import 'package:blckur/exceptions.dart';

import 'package:blckur/collection.dart' as collec;
import 'package:blckur/models/test.dart' as tst;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;
import 'dart:math' as math;

@Injectable()
class Tests extends collec.Collection {
  Type model = tst.Test;

  Tests(ng.Http http) : super(http);
}
