import React from 'react';
import {shallow} from 'enzyme';
import DateInput from './DateInput';
import {Icon, IconButton, TextField} from '@material-ui/core';

const mockEvent = (value: string) => ({
    currentTarget: {value},
});

describe('DateField', () => {
    const props = {
        initialValue: '2020-12-25',
        onDateChange: jest.fn(),
        onDateError: jest.fn(),
    };
    it('displays a text field with a calendar button', () => {
        const component = shallow(<DateInput {...props} />);

        expect(component.find(TextField)).toHaveProp('value', props.initialValue);
        expect(component.find(TextField)).toHaveProp('InputProps',
            expect.objectContaining({
                endAdornment: <IconButton size='small'><Icon>today</Icon></IconButton>,
                autoFocus: true,
            }),
        );
    });
    it('only allows digits and dashes', () => {
        const component = shallow(<DateInput {...props} />);

        component.find(TextField).simulate('change', mockEvent('abc_123-456!'));

        expect(component.find(TextField)).toHaveProp('value', '123-456');
        expect(component.find(TextField)).toHaveClassName('error');
        expect(props.onDateChange).not.toBeCalled();
        expect(props.onDateError).toBeCalledWith('123-456');
    });
    it('accepts single digit month', () => {
        const component = shallow(<DateInput {...props} />);

        component.find(TextField).simulate('change', mockEvent('2020-1-25'));

        expect(component.find(TextField)).toHaveProp('value', '2020-1-25');
        expect(component.find(TextField)).not.toHaveClassName('error');
        expect(props.onDateChange).toBeCalledWith(new Date('2020-01-25'), '2020-01-25');
        expect(props.onDateError).not.toBeCalled();
    });
    it('accepts single digit date', () => {
        const component = shallow(<DateInput {...props} />);

        component.find(TextField).simulate('change', mockEvent('2020-12-2'));

        expect(component.find(TextField)).toHaveProp('value', '2020-12-2');
        expect(props.onDateChange).toBeCalledWith(new Date('2020-12-02'), '2020-12-02');
        expect(props.onDateError).not.toBeCalled();
    });
    it('does not call onDateChange for equivalent value', () => {
        const component = shallow(<DateInput {...props} initialValue='2020-12-02'/>);

        component.find(TextField).simulate('change', mockEvent('2020-12-2'));

        expect(component.find(TextField)).toHaveProp('value', '2020-12-2');
        expect(props.onDateChange).not.toBeCalled();
        expect(props.onDateError).not.toBeCalled();
    });
    describe('onDateError', () => {
        it('is optional', () => {
            const component = shallow(<DateInput {...props} onDateError={undefined} />);

            component.find(TextField).simulate('change', mockEvent('abc_123-456!'));

            expect(component.find(TextField)).toHaveProp('value', '123-456');
        });
    });
});
