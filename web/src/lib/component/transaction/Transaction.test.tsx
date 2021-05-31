import React from 'react';
import {shallow} from 'enzyme';
import Transaction from './Transaction';
import {newTxModel} from 'test/transactionFactory';
import Payee from './Payee';
import Security from './Security';
import Memo from './Memo';
import TxDetail from './TxDetail';
import {newDetail} from 'test/detailFactory';
import DateInput from '../DateInput';
import IconInput from '../IconInput';
import PayeeInput from './PayeeInput';
import SecurityInput from './SecurityInput';
import {Checkbox} from '@material-ui/core';

describe('Transaction', () => {
    const props = {
        tx: newTxModel({referenceNumber: 'check #', memo: 'tx description', details: [newDetail({amount: 123.45})], balance: 345.67}),
        selected: false,
        fieldIndex: 0,
        showSecurity: true,
        setField: jest.fn(),
    };
    const inputProps = {
        variant: 'outlined',
        color: 'primary',
        size: 'small',
    };

    it('displays transaction properties', () => {
        const component = shallow(<Transaction {...props} />);

        expect(component.find('span.date')).toHaveText(props.tx.date);
        expect(component.find('span[data-type="ref"]')).toHaveText(props.tx.referenceNumber!);
        expect(component.find(Memo)).toHaveProp('text', props.tx.memo!);
        expect(component.find(Payee)).toHaveProp('transaction', props.tx);
        expect(component.find(Security)).toHaveProp('transaction', props.tx);
        expect(component.find(TxDetail)).toHaveProp({detail: props.tx.details[0], editField: false, showSecurity: true});
        expect(component.find(Checkbox)).toHaveProp('checked', false);
        expect(component.find('span.number').at(0)).toHaveText(`$${props.tx.subtotal}`);
        expect(component.find('span.number').at(1)).toHaveText(`$${props.tx.balance}`);
    });
    describe('keyboard nav', () => {
        const event = {
            stopPropagation: jest.fn(),
            preventDefault: jest.fn(),
            key: 'Tab',
        } as unknown as React.KeyboardEvent<HTMLDivElement>;

        it('moves to next field on Tab', () => {
            jest.spyOn(props.tx, 'nextField').mockReturnValue(2);
            const component = shallow(<Transaction {...props} selected={true} />);
            const onKeyDown = component.find(DateInput).prop('onKeyDown')!;

            onKeyDown(event);

            expect(props.tx.nextField).toBeCalledWith(1, true);
            expect(props.setField).toBeCalledWith(2);
            expect(event.preventDefault).toBeCalled();
            expect(event.stopPropagation).toBeCalled();
        });
        it('moves to previous field on shift+Tab', () => {
            jest.spyOn(props.tx, 'nextField').mockReturnValue(2);
            const component = shallow(<Transaction {...props} selected={true} />);
            const onKeyDown = component.find(DateInput).prop('onKeyDown')!;

            onKeyDown({...event, shiftKey: true});

            expect(props.tx.nextField).toBeCalledWith(-1, true);
            expect(props.setField).toBeCalledWith(2);
            expect(event.preventDefault).toBeCalled();
            expect(event.stopPropagation).toBeCalled();
        });
        it('ignores ctrl+Tab', () => {
            jest.spyOn(props.tx, 'nextField').mockReturnValue(2);
            const component = shallow(<Transaction {...props} selected={true} />);
            const onKeyDown = component.find(DateInput).prop('onKeyDown')!;

            onKeyDown({...event, ctrlKey: true});

            expect(props.tx.nextField).not.toBeCalled();
            expect(props.setField).not.toBeCalledWith();
            expect(event.preventDefault).not.toBeCalled();
            expect(event.stopPropagation).not.toBeCalled();
        });
        it('ignores alt+Tab', () => {
            jest.spyOn(props.tx, 'nextField').mockReturnValue(2);
            const component = shallow(<Transaction {...props} selected={true} />);
            const onKeyDown = component.find(DateInput).prop('onKeyDown')!;

            onKeyDown({...event, altKey: true});

            expect(props.tx.nextField).not.toBeCalled();
            expect(props.setField).not.toBeCalledWith();
            expect(event.preventDefault).not.toBeCalled();
            expect(event.stopPropagation).not.toBeCalled();
        });
    });
    describe('date input', () => {
        it('shows transaction date', () => {
            const component = shallow(<Transaction {...props} selected={true} />);

            expect(component.find(DateInput)).toHaveProp({
                ...inputProps,
                initialValue: props.tx.date,
            });
        });
        it('updates transaction on change', () => {
            const component = shallow(<Transaction {...props} selected={true} />);

            component.find(DateInput).simulate('dateChange', new Date('2020-12-25'), '2020-12-25');

            expect(props.tx.date).toEqual('2020-12-25');
        });
        it('shows blank for new transaction', () => {
            const component = shallow(<Transaction {...props} selected={true} tx={newTxModel({date: undefined})} />);

            expect(component.find(DateInput)).toHaveProp('initialValue', '');
        });
    });
    describe('ref input', () => {
        it('shows transaction reference number', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={1} />);

            expect(component.find(IconInput)).toHaveProp({
                ...inputProps,
                icon: 'tag',
                value: props.tx.referenceNumber,
            });
        });
        it('updates transaction on change', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={1} />);

            component.find(IconInput).simulate('change', {currentTarget: {value: 'xyz'}});

            expect(props.tx.referenceNumber).toEqual('xyz');
        });
        it('shows blank for new transaction', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={1} tx={newTxModel()} />);

            expect(component.find(IconInput)).toHaveProp('value', '');
        });
    });
    describe('payee input', () => {
        it('shows transaction payee', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={2} />);

            expect(component.find(PayeeInput)).toHaveProp({
                ...inputProps,
                transaction: props.tx,
            });
        });
    });
    describe('security input', () => {
        it('shows transaction security', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={3} />);

            expect(component.find(SecurityInput)).toHaveProp({
                ...inputProps,
                transaction: props.tx,
            });
        });
    });
    describe('description input', () => {
        it('shows transaction security', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={4} />);

            expect(component.find(IconInput)).toHaveProp({
                ...inputProps,
                icon: 'notes',
                value: props.tx.memo,
            });
        });
        it('updates transaction on change', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={4} />);

            component.find(IconInput).simulate('change', {currentTarget: {value: 'xyz'}});

            expect(props.tx.memo).toEqual('xyz');
        });
        it('shows blank for no memo', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={4} tx={newTxModel()} />);

            expect(component.find(IconInput)).toHaveProp('value', '');
        });
    });
    describe('cleared checkbox', () => {
        it('updates transaction on change', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={4} />);

            component.find(Checkbox).simulate('change', {}, true);

            expect(props.tx.cleared).toBe(true);
        });
    });
    describe('detail inputs', () => {
        const detailProps = {...inputProps, detail: props.tx.details[0]};
        it('shows amount input', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={5} />);

            expect(component.find(TxDetail)).toHaveProp({...detailProps, editField: 'amount'});
        });
        it('shows category input', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={6} />);

            expect(component.find(TxDetail)).toHaveProp({...detailProps, editField: 'category'});
        });
        it('shows group input', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={7} />);

            expect(component.find(TxDetail)).toHaveProp({...detailProps, editField: 'group'});
        });
        it('shows shares input', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={8} />);

            expect(component.find(TxDetail)).toHaveProp({...detailProps, editField: 'shares'});
        });
        it('shows memo input', () => {
            const component = shallow(<Transaction {...props} selected={true} fieldIndex={9} />);

            expect(component.find(TxDetail)).toHaveProp({...detailProps, editField: 'memo'});
        });
    });
});
