library recaptcha_comp;

import 'package:blckur/alert.dart' as alrt;

import 'package:angular/angular.dart' show Component;
import 'package:angular/angular.dart' as ng;
import 'dart:html' as dom;
import 'dart:js' as js;
import 'dart:async' as async;

@Component(
  selector: 'x-recaptcha',
  template: '<div vertical layout center></div>'
)
class RecaptchaComp implements ng.ShadowRootAware {
  dom.ShadowRoot root;
  dom.ScriptElement apiElem;
  DateTime start;

  void _init() {
    new async.Timer(const Duration(milliseconds: 5), () {
      var api = js.context['grecaptcha'];

      if (new DateTime.now().difference(this.start) > const Duration(
          seconds: 10)) {
        new alrt.Alert('Unable to load reCaptcha', () {
          this.loadRecaptch();
        });
        return;
      }

      if (api == null) {
        this._init();
        return;
      }

      var elem = new dom.DivElement();
      this.root.querySelector('div').append(elem);

      js.context['grecaptcha'].callMethod('render', [
        elem,
        new js.JsObject.jsify({
          'sitekey': '6LeVEgMTAAAAAAdu7KTYPUZG5deiVER-84TT1OXf',
          'theme': 'dark',
          'type': 'image',
          'callback': (resp) {
            print(resp);
          },
        }),
      ]);

    });
  }

  void loadRecaptch() {
    this.start = new DateTime.now();

    if (this.apiElem != null) {
      this.apiElem.remove();
    }

    this.apiElem = new dom.ScriptElement()
      ..type = 'application/javascript'
      ..src = '//www.google.com/recaptcha/api.js';

    this.root.querySelector('div').append(this.apiElem);

    this._init();
  }

  void onShadowRoot(dom.ShadowRoot root) {
    this.root = root;
    this.loadRecaptch();
  }
}
