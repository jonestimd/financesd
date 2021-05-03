import React from 'react';
import {shallow} from 'enzyme';
import Memo from './Memo';
import {Icon} from '@material-ui/core';

describe('Memo', () => {
    it('returns null if no memo', () => {
        expect(shallow(<Memo />)).toBeEmptyRender();
    });
    it('displays memo', () => {
        const memo = 'the memo';

        const component = shallow(<Memo text={memo} />);

        expect(component).toHaveProp('data-type', 'description');
        expect(component.find(Icon)).toHaveText('notes');
        expect(component.childAt(1)).toHaveText(memo);
    });
});
