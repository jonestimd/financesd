import React from 'react';
import {shallow} from 'enzyme';
import SecurityInput from './SecurityInput';
import Autocomplete, {AutocompleteRenderInputParams} from '@material-ui/lab/Autocomplete';
import {RootStore} from 'src/lib/store/RootStore';
import {newSecurityModel} from 'src/test/securityFactory';
import {Icon} from '@material-ui/core';
import {newTxModel} from 'src/test/transactionFactory';

describe('SecurityInput', () => {
    const {securityStore} = new RootStore();
    const securities = [newSecurityModel({name: 'Some Stock', symbol: 'STK'}), newSecurityModel({name: 'A Mutual Fund'})];
    securities.forEach((s) => securityStore['securitiesById'].set(s.id, s));

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({securityStore});
    });
    it('displays list of accounts and categories', () => {
        const transaction = newTxModel();

        const component = shallow(<SecurityInput transaction={transaction} />);

        expect(component.find(Autocomplete)).toHaveProp('options', securityStore.securities);
        expect(component.find(Autocomplete)).toHaveProp('value', null);
    });
    it('uses displayName for option labels', () => {
        const transaction = newTxModel();
        const component = shallow(<SecurityInput transaction={transaction} />);

        const getOptionLabel = component.find(Autocomplete).prop('getOptionLabel')!;

        expect(getOptionLabel(securities[0])).toEqual(securities[0].displayName);
    });
    it('displays options containing text, ignoring case', () => {
        const transaction = newTxModel();
        const component = shallow(<SecurityInput transaction={transaction} />);

        const filterOptions = component.find(Autocomplete).prop('filterOptions')!;

        expect(filterOptions(securities, {inputValue: 'fund', getOptionLabel: () => ''})).toEqual(securities.slice(1));
    });
    it('displays selected security', () => {
        const transaction = newTxModel({securityId: securities[1].id});

        const component = shallow(<SecurityInput transaction={transaction} />);

        expect(component.find(Autocomplete)).toHaveProp('value', securities[1]);
    });
    it('displays input with request_page icon', () => {
        const transaction = newTxModel();
        const component = shallow(<SecurityInput transaction={transaction} />);
        const renderInput = component.find(Autocomplete).prop('renderInput')!;

        const input = shallow(renderInput({} as AutocompleteRenderInputParams) as React.ReactElement);

        expect(input).toHaveProp('InputProps', expect.objectContaining({
            startAdornment: <Icon>request_page</Icon>,
        }));
    });
    it('updates transaction when selection changes', () => {
        const transaction = newTxModel({securityId: securities[1].id});
        const component = shallow(<SecurityInput transaction={transaction} />);

        component.find(Autocomplete).simulate('change', {}, securities[0]);

        expect(transaction.securityId).toEqual(securities[0].id);
    });
    it('updates transaction when selection is cleared', () => {
        const transaction = newTxModel({securityId: securities[1].id});
        const component = shallow(<SecurityInput transaction={transaction} />);

        component.find(Autocomplete).simulate('change', {}, null);

        expect(transaction.securityId).toBeUndefined();
    });
});
