import TransactionTableModel from '../model/TransactionTableModel';
import {RootStore} from './RootStore';
import * as agent from '../agent';
import {loadingTransactions, query} from './TransactionStore';

describe('TransactionStore', () => {
    const accountId = '1';
    const {transactionStore, categoryStore, messageStore, alertStore} = new RootStore();
    const tableModel = new TransactionTableModel([], categoryStore);

    beforeEach(() => {
        transactionStore['transactionsByAccountId'].clear();
    });
    describe('getTransactionsModel', () => {
        it('returns empty model for unknown account', () => {
            expect(transactionStore.getTransactionsModel('-99')).toBe(TransactionTableModel.EMPTY);
        });
        it('returns table model for account', () => {
            transactionStore['transactionsByAccountId'].set(accountId, tableModel);

            expect(transactionStore.getTransactionsModel(accountId)).toBe(tableModel);
        });
    });
    describe('loadTransactions', () => {
        beforeEach(() => {
            jest.spyOn(messageStore, 'addProgressMessage');
            jest.spyOn(messageStore, 'removeProgressMessage');
            transactionStore['pendingAccounts'].splice(0, transactionStore['pendingAccounts'].length);
        });
        it('does nothing is transactions already loaded', async () => {
            transactionStore['transactionsByAccountId'].set(accountId, tableModel);
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await transactionStore.loadTransactions(accountId);

            expect(agent.graphql).not.toBeCalled();
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
        });
        it('does nothing is transactions already loading', async () => {
            transactionStore['pendingAccounts'].push(accountId);
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await transactionStore.loadTransactions(accountId);

            expect(agent.graphql).not.toBeCalled();
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
        });
        it('loads transactions for account', async () => {
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: {transactions: []}});

            await transactionStore.loadTransactions(accountId);

            expect(transactionStore['pendingAccounts']).toHaveLength(0);
            expect(transactionStore.getTransactionsModel(accountId)).toBeInstanceOf(TransactionTableModel);
            expect(agent.graphql).not.toBeCalledWith(query);
            expect(messageStore.addProgressMessage).toBeCalledWith(loadingTransactions);
            expect(messageStore.removeProgressMessage).toBeCalledWith(loadingTransactions);
        });
        it('logs error from graphql', async () => {
            const error = new Error('API error');
            jest.spyOn(agent, 'graphql').mockRejectedValue(error);
            jest.spyOn(console, 'error').mockImplementation(() => { });
            jest.spyOn(alertStore, 'addAlert').mockReturnValue();

            await transactionStore.loadTransactions(accountId);

            expect(transactionStore['pendingAccounts']).toHaveLength(0);
            expect(transactionStore.getTransactionsModel(accountId)).toBe(TransactionTableModel.EMPTY);
            expect(alertStore.addAlert).toBeCalledWith('error', 'Error loading transactions');
            expect(console.error).toBeCalledWith('error from Loading transactions', error);
        });
    });
});
