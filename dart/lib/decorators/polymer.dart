library polymer_dec;

import 'package:angular/angular.dart' show Decorator;
import 'package:angular/angular.dart' as ng;
import 'package:observe/observe.dart' as obsrv;
import 'package:template_binding/template_binding.dart' as tmpl_bnd;
import 'dart:html' as dom;

typedef dynamic _Nullary();

Map<String, String> POLYMER_ATTR_MAP = const {
  'py-active': 'active',
  'py-animated': 'animated',
  'py-autoclosedisabled': 'autoclosedisabled',
  'py-checked': 'checked',
  'py-closeselector': 'closeselector',
  'py-closedicon': 'closedicon',
  'py-committedvalue': 'committedvalue',
  'py-disabledrag': 'disabledrag',
  'py-disabled': 'disabled',
  'py-duration': 'duration',
  'py-editable': 'editable',
  'py-error': 'error',
  'py-fill': 'fill',
  'py-floatinglabel': 'floatinglabel',
  'py-heading': 'heading',
  'py-hidescrollbutton': 'hidescrollbutton',
  'py-icon': 'icon',
  'py-immediatevalue': 'immediatevalue',
  'py-indeterminate': 'indeterminate',
  'py-initialopacity': 'initialopacity',
  'py-isinvalid': 'isinvalid',
  'py-label': 'label',
  'py-labelvisible': 'labelvisible',
  'py-layered': 'layered',
  'py-max': 'max',
  'py-maxrows': 'maxrows',
  'py-min': 'min',
  'py-mini': 'mini',
  'py-nobar': 'nobar',
  'py-noink': 'noink',
  'py-noslide': 'noslide',
  'py-opacitydecayvelocity': 'opacitydecayvelocity',
  'py-opened': 'opened',
  'py-openedicon': 'openedicon',
  'py-pin': 'pin',
  'py-raised': 'raised',
  'py-recenteringtouch': 'recenteringtouch',
  'py-responsivewidth': 'responsivewidth',
  'py-rows': 'rows',
  'py-scrollable': 'scrollable',
  'py-secondaryprogress': 'secondaryprogress',
  'py-snaps': 'snaps',
  'py-src': 'src',
  'py-swipedisabled': 'swipedisabled',
  'py-target': 'target',
  'py-text': 'text',
  'py-toggles': 'toggles',
  'py-transition': 'transition',
  'py-value': 'value',
  'py-z': 'z',
};

@Decorator(
  selector: '[py-*]',
  exportExpressionAttrs: const [
    'py-active',
    'py-animated',
    'py-autoCloseDisabled',
    'py-checked',
    'py-closeSelector',
    'py-closedIcon',
    'py-committedValue',
    'py-disableDrag',
    'py-disabled',
    'py-duration',
    'py-editable',
    'py-error',
    'py-fill',
    'py-floatingLabel',
    'py-heading',
    'py-hideScrollButton',
    'py-icon',
    'py-immediateValue',
    'py-indeterminate',
    'py-initialOpacity',
    'py-isInvalid',
    'py-label',
    'py-labelVisible',
    'py-layered',
    'py-max',
    'py-maxRows',
    'py-min',
    'py-mini',
    'py-nobar',
    'py-noink',
    'py-noslide',
    'py-opacityDecayVelocity',
    'py-opened',
    'py-openedIcon',
    'py-pin',
    'py-raised',
    'py-recenteringTouch',
    'py-responsiveWidth',
    'py-rows',
    'py-scrollable',
    'py-secondaryProgress',
    'py-snaps',
    'py-src',
    'py-swipeDisabled',
    'py-target',
    'py-text',
    'py-toggles',
    'py-transition',
    'py-value',
    'py-z',
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
      var polymerKey = POLYMER_ATTR_MAP[key];

      if (polymerKey == null) {
        throw new StateError('Unknown polymer attribute ${key}');
      }

      tmpl_bnd.nodeBind(node).bind(polymerKey, bindVal);

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
