library delay_click_dec;

import 'package:blckur/utils/utils.dart' as utils;

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

    // Cache missed clicks on short touch moves on ios
    var start;
    var startPt;
    var endPt;

    element.onTouchStart.listen((dom.TouchEvent evt) {
      if (evt.touches.length == 0) {
        return;
      }

      start = new DateTime.now();
      startPt = evt.touches[0].page;
      endPt = null;
    });

    element.onTouchLeave.listen((dom.TouchEvent evt) {
      start = null;
    });

    element.onTouchMove.listen((dom.TouchEvent evt) {
      if (evt.touches.length == 0) {
        return;
      }
      endPt = evt.touches[0].page;
    });

    element.onTouchEnd.listen((evt) {
      if (start == null || endPt == null) {
        return;
      }

      var diff = new DateTime.now().difference(start);
      var dist = utils.distance(startPt.x, startPt.y, endPt.x, endPt.y);

      if (dist < 30 && diff < const Duration(milliseconds: 250)) {
        new async.Timer(const Duration(milliseconds: 150), () {
          this.callback({r'$event': evt});
        });
      }
    });
  }
}
