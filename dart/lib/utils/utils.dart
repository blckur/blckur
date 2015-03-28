library utils;

export 'package:blckur/utils/loading.dart' show Loading;
export 'package:blckur/utils/modeswitch.dart' show ModeSwitch;

import 'dart:html' as dom;
import 'dart:math' as math;
import 'dart:js' as js;

String uuid() {
  var rand = new math.Random();
  var id = '';

  for (var i = 0; i < 8; i++) {
    id += ((1 + rand.nextDouble()) * 0x10000).floor().toRadixString(
      16).substring(1);
  }

  return id;
}

String getDomain() {
  return '${dom.window.location.protocol}//${dom.window.location.host}';
}

void printColor(dynamic obj, String color) {
  js.context['console'].callMethod('log', [
    '%c$obj', 'color: $color']);
}
