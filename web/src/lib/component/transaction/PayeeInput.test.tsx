import React from 'react';
import {shallow} from 'enzyme';
import PayeeInput from './PayeeInput';
import Autocomplete, {AutocompleteRenderInputParams} from '@material-ui/lab/Autocomplete';
import {RootStore} from 'src/lib/store/RootStore';
import {newPayee} from 'src/test/payeeFactory';
import {Icon} from '@material-ui/core';
import {newTxModel} from 'src/test/transactionFactory';

describe('PayeeInput', () => {
    const {payeeStore} = new RootStore();
    const payees = [newPayee({name: 'Somebody'}), newPayee({name: 'anOther perSON'})];
    payees.forEach((p) => payeeStore['payeesById'].set(p.id, p));

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({payeeStore});
    });
    it('displays list of accounts and categories', () => {
        const transaction = newTxModel();

        const component = shallow(<PayeeInput transaction={transaction} />);

        expect(component.find(Autocomplete)).toHaveProp('options', payeeStore.payees);
        expect(component.find(Autocomplete)).toHaveProp('value', null);
    });
    it('uses name for option labels', () => {
        const transaction = newTxModel();
        const component = shallow(<PayeeInput transaction={transaction} />);

        const getOptionLabel = component.find(Autocomplete).prop('getOptionLabel')!;

        expect(getOptionLabel(payees[0])).toEqual(payees[0].name);
    });
    it('displays options containing text, ignoring case', () => {
        const transaction = newTxModel();
        const component = shallow(<PayeeInput transaction={transaction} />);

        const filterOptions = component.find(Autocomplete).prop('filterOptions')!;

        expect(filterOptions(payees, {inputValue: 'BODY', getOptionLabel: () => ''})).toEqual(payees.slice(0,1));
    });
    it('displays selected payee', () => {
        const transaction = newTxModel({payeeId: payees[1].id});

        const component = shallow(<PayeeInput transaction={transaction} />);

        expect(component.find(Autocomplete)).toHaveProp('value', payees[1]);
    });
    it('displays input with person icon', () => {
        const transaction = newTxModel();
        const component = shallow(<PayeeInput transaction={transaction} />);
        const renderInput = component.find(Autocomplete).prop('renderInput')!;

        const input = shallow(renderInput({} as AutocompleteRenderInputParams) as React.ReactElement);

        expect(input).toHaveProp('InputProps', expect.objectContaining({
            startAdornment: <Icon>person</Icon>,
        }));
    });
});
