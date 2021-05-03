import React from 'react';
import {shallow} from 'enzyme';
import {Icon, TextField} from '@material-ui/core';
import IconInput from './IconInput';

describe('IconInput', () => {
    it('displays a text field with a leading icon', () => {
        const component = shallow(<IconInput icon='notes' />);

        expect(component.find(TextField)).toHaveProp('InputProps',
            expect.objectContaining({
                startAdornment: <Icon>notes</Icon>,
                autoFocus: true,
            }),
        );
    });
});
