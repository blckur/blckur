library polymer_dec;

import 'package:angular/angular.dart' show Decorator;
import 'package:angular/angular.dart' as ng;
import 'package:observe/observe.dart' as obsrv;
import 'package:template_binding/template_binding.dart' as tmpl_bnd;
import 'dart:html' as dom;

typedef dynamic _Nullary();

@Decorator(
  selector: '[py-*]',
  exportExpressionAttrs: const [
    'py-mini',
    'py-raised',
    'py-closeSelector',
    'py-pin',
    'py-text',
    'py-autoCloseDisabled',
    'py-autoValidate',
    'py-disabled',
    'py-src',
    'py-duration',
    'py-alt',
    'py-dragging',
    'py-labelVisible',
    'py-fill',
    'py-rows',
    'py-checked',
    'py-opened',
    'py-responsiveWidth',
    'py-indeterminate',
    'py-label',
    'py-immediateValue',
    'py-hideScrollButton',
    'py-icon',
    'py-secondaryProgress',
    'py-committedValue',
    'py-snaps',
    'py-scrollable',
    'py-toggles',
    'py-editable',
    'py-recenteringTouch',
    'py-closedIcon',
    'py-initialOpacity',
    'py-openedIcon',
    'py-noslide',
    'py-active',
    'py-noink',
    'py-swipeDisabled',
    'py-nobar',
    'py-disableDrag',
    'py-z',
    'py-target',
    'py-isInvalid',
    'py-transition',
    'py-value',
    'py-layered',
    'py-maxRows',
    'py-error',
    'py-floatingLabel',
    'py-opacityDecayVelocity',
    'py-animated',
    'py-heading',
  ]
)
class PolymerDec {
  PolymerDec(dom.Node node, ng.Parser parser, ng.Scope scope) {
    dom.Element element = node;

    element.attributes.forEach((key, value) {
      if (!key.startsWith('py-')) {
        return;
      }
      var bindVal = new ValueBindable();
      var expression = parser(value);

      tmpl_bnd.nodeBind(node).bind(key.substring(3), bindVal);

      if (expression.isAssignable) {
        bindVal.onChange = (value) {
          expression.assign(scope.context, value);
        };
      }

      scope.watch(value, bindVal.update);
    });
  }
}

class ValueBindable implements obsrv.Bindable {
  Function callback;
  Function onChange;

  var _value;
  void set value(dynamic value) {
    this._value = value;
    if (this.onChange != null) {
      this.onChange(this._value);
    }
  }
  dynamic get value {
    return this._value;
  }

  void close() {
    this.callback = null;
  }

  dynamic open(Function callback) {
    this.callback = callback;
    return this._value;
  }

  void update(dynamic newValue, dynamic oldValue) {
    this._value = newValue;
    if (this.callback != null) {
      if (this.callback is _Nullary) {
        this.callback();
      }
      else {
        this.callback(this._value);
      }
    }
  }

  void deliver() {}
}
