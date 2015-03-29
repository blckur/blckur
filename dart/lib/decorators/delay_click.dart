library delay_click_dec;

import 'package:angular/angular.dart' show Decorator, NgCallback;
import 'dart:html' as dom;
import 'dart:async' as async;

@Decorator(
  selector: '[delay-click]'
)
class DelayClickDec {
  @NgCallback('delay-click')
  Function callback;

  DelayClickDec(dom.Element element) {
    element.onClick.listen((evt) {
      new async.Timer(const Duration(milliseconds: 185), () {
        this.callback({r'$event': evt});
      });
    });
  }
}
