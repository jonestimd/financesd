import React from 'react';
import {shallow} from 'enzyme';
import DateInput from './DateInput';
import {Icon, IconButton, TextField} from '@material-ui/core';

describe('DateField', () => {
    it('displays a text field with a calendar button', () => {
        const component = shallow(<DateInput />);

        expect(component.find(TextField)).toHaveProp('InputProps',
            expect.objectContaining({
                endAdornment: <IconButton size='small'><Icon>today</Icon></IconButton>,
                autoFocus: true,
            }),
        );
    });
});
