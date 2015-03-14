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
    element.onClick.listen((_) {
      new async.Timer(const Duration(milliseconds: 150), () {
        this.callback();
      });
    });
  }
}
