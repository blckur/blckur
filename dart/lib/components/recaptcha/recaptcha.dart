library recaptcha_comp;

import 'package:blckur/utils/utils.dart' as utils;

import 'package:angular/angular.dart' show Component;
import 'package:angular/angular.dart' as ng;
import 'dart:html' as dom;
import 'dart:js' as js;

@Component(
  selector: 'x-recaptcha',
  template: '<div vertical layout center></div>'
)
class RecaptchaComp implements ng.ShadowRootAware {
  dom.ShadowRoot root;
  dom.ScriptElement apiElem;
  DateTime start;

  void _init() {
    var api = js.context['grecaptcha'];

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
  }

  void loadRecaptch() {
    this.start = new DateTime.now();

    if (this.apiElem != null) {
      this.apiElem.remove();
    }

    var callback = 'cb' + utils.uuid();

    js.context[callback] = () {
      this._init();
    };

    this.apiElem = new dom.ScriptElement()
      ..type = 'text/javascript'
      ..async = true
      ..defer = true
      ..src = '//www.google.com/recaptcha/api.js?'
        'onload=${callback}&render=explicit';

    this.root.querySelector('div').append(this.apiElem);
  }

  void onShadowRoot(dom.ShadowRoot root) {
    this.root = root;
    this.loadRecaptch();
  }
}
