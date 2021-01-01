import React from 'react';
import {shallow} from 'enzyme';
import Memo from './Memo';

describe('Memo', () => {
    it('returns null if no memo', () => {
        expect(shallow(<Memo text={null} />)).toBeEmptyRender();
    });
    it('displays memo', () => {
        const memo = 'the memo';

        const component = shallow(<Memo text={memo} />);

        expect(component.find('i.material-icons.md-18')).toHaveText('notes');
        expect(component.find('.memo.chip').childAt(1)).toHaveText(memo);
    });
});
