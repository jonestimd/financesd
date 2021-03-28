import React, {ChangeEvent} from 'react';
import {shallow} from 'enzyme';
import {History} from 'history';
import {RootStore} from 'src/lib/store/RootStore';
import {newAccountModel, newCompanyModel} from 'src/test/accountFactory';
import TransactionsPage from './TransactionsPage';
import TransactionList from './TransactionList';
import TransactionTable from './TransactionTable';
import ToggleButtonGroup from '@material-ui/lab/ToggleButtonGroup';
import ToggleButton from '@material-ui/lab/ToggleButton';
import TopAppBar from '../TopAppBar';
import {mockUseEffect} from 'src/test/mockHooks';
import Autocomplete, {AutocompleteRenderInputParams} from '@material-ui/lab/Autocomplete';
import {AccountModel} from 'src/lib/model/account/AccountModel';
import {TextField} from '@material-ui/core';
import settingsStore from 'src/lib/store/settingsStore';

const accountId = '123';
const history = {push: jest.fn()} as unknown as History;

jest.mock('react-router', () => ({
    ...jest.requireActual('react-router'),
    useHistory: () => history,
// eslint-disable-next-line @typescript-eslint/ban-types
} as object));

describe('TransactionsPage', () => {
    const {accountStore, transactionStore} = new RootStore();
    const account = newAccountModel({}, newCompanyModel());
    const props = {
        match: {params: {accountId}},
    };

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({accountStore, transactionStore});
        jest.spyOn(transactionStore, 'loadTransactions').mockResolvedValue(true);
        jest.spyOn(accountStore, 'getAccount').mockReturnValue(account);
        mockUseEffect();
    });
    it('loads account transactions', () => {
        shallow(<TransactionsPage {...props} />);

        expect(transactionStore.loadTransactions).toBeCalledWith(parseInt(accountId));
        expect(accountStore.getAccount).toBeCalledWith(parseInt(accountId));
    });
    describe('account input', () => {
        it('displays account name', () => {
            const component = shallow(<TransactionsPage {...props} />);

            expect(component.find(TopAppBar).find(Autocomplete)).toHaveProp('value', account);
        });
        it('displays blank account name', () => {
            jest.spyOn(accountStore, 'getAccount').mockReturnValue(undefined);

            const component = shallow(<TransactionsPage {...props} />);

            expect(component.find(TopAppBar).find(Autocomplete)).toHaveProp('value', null);
        });
        it('groups accounts by company', () => {
            type GroupByProp = (a: AccountModel) => string;
            const groupBy = shallow(<TransactionsPage {...props} />).find(Autocomplete).prop<GroupByProp>('groupBy');

            expect(groupBy(account)).toEqual(account.company!.name);
            expect(groupBy(newAccountModel())).toEqual('');
        });
        it('displays company and account names', () => {
            type GetOptionLabelProp = (a: AccountModel) => string;
            const getOptionLabel = shallow(<TransactionsPage {...props} />).find(Autocomplete).prop<GetOptionLabelProp>('getOptionLabel');

            expect(getOptionLabel(account)).toEqual(account.displayName);
        });
        it('displays account name in list', () => {
            type RenderOptionProp = (a: AccountModel) => string;
            const renderOption = shallow(<TransactionsPage {...props} />).find(Autocomplete).prop<RenderOptionProp>('renderOption');

            expect(renderOption(account)).toEqual(account.name);
        });
        it('uses TextField for input', () => {
            const renderInput = shallow(<TransactionsPage {...props}/>).find(Autocomplete).prop('renderInput');
            const params = {id: 'x', disabled: false} as AutocompleteRenderInputParams;

            expect(renderInput(params)).toEqual(<TextField {...params} variant='outlined'/>);
        });
        it('displays transactions for selected account', () => {
            const component = shallow(<TransactionsPage {...props}/>);
            const onChange = component.find(Autocomplete).prop('onChange')!;

            onChange({} as ChangeEvent, account, 'select-option');

            expect(history.push).toBeCalledWith(`/finances/account/${account.id}`);
        });
    });
    describe('mode', () => {
        it('displays list and table mode buttons', () => {
            const component = shallow(<TransactionsPage {...props} />);

            const buttons = component.find(ToggleButtonGroup).find(ToggleButton);
            expect(buttons).toHaveLength(2);
            expect(buttons.map((b) => b.prop<string>('value'))).toEqual(['list', 'table']);
        });
        it('uses settings', () => {
            settingsStore.transactionsView = 'list';

            const component = shallow(<TransactionsPage {...props} />);

            expect(component.find(TransactionList)).toHaveProp('accountId', parseInt(accountId));
            expect(component.find(TransactionTable)).not.toExist();
        });
        it('shows table when table button clicked', () => {
            settingsStore.transactionsView = 'list';
            const component = shallow(<TransactionsPage {...props} />);

            component.find(ToggleButtonGroup).simulate('change', {}, 'table');
            component.rerender({...props, x: 'force rerender'});

            expect(component.find(TransactionList)).not.toExist();
            expect(component.find(TransactionTable)).toHaveProp('accountId', parseInt(accountId));
            expect(settingsStore.transactionsView).toEqual('table');
        });
    });
});
