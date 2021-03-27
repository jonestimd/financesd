import React from 'react';
import {shallow} from 'enzyme';
import ChildMenu from './ChildMenu';
import {MenuItem} from '@material-ui/core';

describe('ChildMenu', () => {
    const onBack = jest.fn();
    it('displays back button and children', () => {
        const component = shallow(<ChildMenu onBack={onBack}><div id='children'/></ChildMenu>);

        expect(component.find(MenuItem).text()).toEqual(expect.stringMatching(/Back$/));
        component.find(MenuItem).simulate('click');
        expect(onBack).toBeCalled();
        expect(component.find('#children')).toExist();
    });
    it('displays title', () => {
        const name = 'menu name';

        const component = shallow(<ChildMenu onBack={onBack} name={name}><div/></ChildMenu>);

        const title = component.find(MenuItem).at(1);
        expect(title).toHaveText(name);
        expect(title).toHaveProps({disabled: true, divider: true});
    });
});
