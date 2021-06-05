import React from 'react';
import {shallow} from 'enzyme';
import NumberInput from './NumberInput';
import {TextField} from '@material-ui/core';

describe('NumberInput', () => {
    const props = {
        value: '42',
        precision: 2,
        onChange: jest.fn(),
    };
    it('sets start adornment', () => {
        const component = shallow(<NumberInput {...props} startAdornment={<span>$</span>} />);

        expect(component.find(TextField)).toHaveProp('InputProps', expect.objectContaining({autoFocus: true, startAdornment: <span>$</span>}));
    });
    it('calls onChange for valid input', () => {
        const component = shallow(<NumberInput {...props} />);

        component.find(TextField).simulate('change', {currentTarget: {value: '42.12'}});

        expect(props.onChange).toBeCalledWith('42.12');
        expect(component.find(TextField)).toHaveProp('value', '42.12');
        expect(component.find(TextField)).not.toHaveClassName('error');
    });
    it('ignores onChange for invalid input', () => {
        const component = shallow(<NumberInput {...props} />);

        component.find(TextField).simulate('change', {currentTarget: {value: '42.123'}});

        expect(props.onChange).not.toBeCalled();
        expect(component.find(TextField)).toHaveProp('value', '42');
        expect(component.find(TextField)).not.toHaveClassName('error');
    });
    it('shows error for empty value', () => {
        const component = shallow(<NumberInput {...props} />);

        component.find(TextField).simulate('change', {currentTarget: {value: ''}});

        expect(props.onChange).toBeCalledWith('');
        expect(component.find(TextField)).toHaveProp('value', '');
        expect(component.find(TextField)).toHaveClassName('error');
    });
    it('shows error for invalid value', () => {
        const component = shallow(<NumberInput {...props} />);

        component.find(TextField).simulate('change', {currentTarget: {value: '.'}});

        expect(props.onChange).toBeCalledWith('.');
        expect(component.find(TextField)).toHaveProp('value', '.');
        expect(component.find(TextField)).toHaveClassName('error');
    });
});
