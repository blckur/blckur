library enter_aware;

import 'package:angular/angular.dart' as ng;
import 'dart:html' as dom;

class EnterAware implements ng.ShadowRootAware {
  void onShadowRoot(dom.ShadowRoot root) {
    root.host.onKeyPress.listen((evt) {
      if (evt.keyCode == 13) {
        evt.stopPropagation();
        this.onEnter();
      }
    });
  }

  void onEnter() {}
}
