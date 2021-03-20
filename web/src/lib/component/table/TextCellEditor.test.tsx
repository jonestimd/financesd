import React from 'react';
import {shallow} from 'enzyme';
import TextCellEditor from './TextCellEditor';

describe('TextCellEditor', () => {
    const value = 'starting value';
    const onCommit = jest.fn();
    const props = {value, onCommit};

    it('restores value on escape', () => {
        const stopPropagation = jest.fn();
        const component = shallow(<TextCellEditor {...props} />);

        component.simulate('change', {target: {value: 'new value'}});
        component.simulate('keydown', {key: 'Escape', stopPropagation});

        expect(stopPropagation).toBeCalledTimes(1);
        expect(onCommit).toBeCalledWith(value);
    });
    it('commits value on enter', () => {
        const stopPropagation = jest.fn();
        const component = shallow(<TextCellEditor {...props} />);

        component.simulate('change', {target: {value: 'new value'}});
        component.simulate('keydown', {key: 'Enter', stopPropagation});

        expect(stopPropagation).not.toBeCalled();
        expect(onCommit).toBeCalledWith('new value');
    });
    it('commits value on tab', () => {
        const stopPropagation = jest.fn();
        const component = shallow(<TextCellEditor {...props} />);

        component.simulate('change', {target: {value: 'new value'}});
        component.simulate('keydown', {key: 'Tab', stopPropagation});

        expect(stopPropagation).not.toBeCalled();
        expect(onCommit).toBeCalledWith('new value');
    });
    it('does not commit on text change', () => {
        const stopPropagation = jest.fn();
        const component = shallow(<TextCellEditor {...props} />);

        component.simulate('change', {target: {value: 'new value'}});
        component.simulate('keydown', {key: 'x', stopPropagation});

        expect(stopPropagation).not.toBeCalled();
        expect(onCommit).not.toBeCalled();
    });
    it('commits value on blur', () => {
        const component = shallow(<TextCellEditor {...props} />);

        component.simulate('change', {target: {value: 'new value'}});
        component.simulate('blur');

        expect(onCommit).toBeCalled();
    });
});
