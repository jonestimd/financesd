import React from 'react';
import {shallow} from 'enzyme';
import {RootStore} from 'src/lib/store/RootStore';
import {newAccount} from 'src/test/accountFactory';
import TransactionsPage from './TransactionsPage';
import TransactionList from './TransactionList';
import TransactionTable from './TransactionTable';
import ToggleButtonGroup from '@material-ui/lab/ToggleButtonGroup';
import ToggleButton from '@material-ui/lab/ToggleButton';
import TopAppBar from '../TopAppBar';

const accountId = '123';

describe('TransactionsPage', () => {
    const {accountStore, transactionStore} = new RootStore();
    const account = newAccount();
    const props = {
        match: {params: {accountId}},
    };

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({accountStore, transactionStore});
        jest.spyOn(transactionStore, 'loadTransactions').mockResolvedValue();
        jest.spyOn(accountStore, 'getAccount').mockReturnValue(account);
        jest.spyOn(React, 'useEffect').mockImplementation((cb) => cb());
    });
    it('loads account transactions', () => {
        shallow(<TransactionsPage {...props} />);

        expect(transactionStore.loadTransactions).toBeCalledWith(accountId);
        expect(accountStore.getAccount).toBeCalledWith(accountId);
    });
    it('displays account name in app bar', () => {
        const component = shallow(<TransactionsPage {...props} />);

        expect(component.find(TopAppBar)).toHaveProp('title', account.name);
    });
    describe('mode', () => {
        it('displays list and table mode buttons', () => {
            const component = shallow(<TransactionsPage {...props} />);

            const buttons = component.find(ToggleButtonGroup).find(ToggleButton);
            expect(buttons).toHaveLength(2);
            expect(buttons.map((b) => b.prop<string>('value'))).toEqual(['list', 'table']);
        });
        it('defaults to list', () => {
            const component = shallow(<TransactionsPage {...props} />);

            expect(component.find(TransactionList)).toHaveProp('accountId', accountId);
            expect(component.find(TransactionTable)).not.toExist();
        });
        it('shows table when table button clicked', () => {
            const component = shallow(<TransactionsPage {...props} />);

            component.find(ToggleButtonGroup).simulate('change', {}, 'table');
            component.rerender({...props, x: 'force rerender'});

            expect(component.find(TransactionList)).not.toExist();
            expect(component.find(TransactionTable)).toHaveProp('accountId', accountId);
        });
    });
});
