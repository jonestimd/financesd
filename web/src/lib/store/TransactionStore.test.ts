import TransactionTableModel from '../model/TransactionTableModel';
import {RootStore} from './RootStore';
import * as agent from '../agent';
import {loadingTransactions, query, savingTransactions, updateTxMutation} from './TransactionStore';
import {newTx} from 'test/transactionFactory';

describe('TransactionStore', () => {
    const accountId = 1;
    const {transactionStore, accountStore, categoryStore, messageStore, alertStore} = new RootStore();
    const tableModel = new TransactionTableModel([], accountStore, categoryStore);

    beforeEach(() => {
        transactionStore['transactionsByAccountId'].clear();
    });
    describe('getTransactionsModel', () => {
        it('returns empty model for unknown account', () => {
            expect(transactionStore.getTransactionsModel(-99)).toBe(TransactionTableModel.EMPTY);
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
    describe('updateTransactions', () => {
        it('ignores unknown account id', async () => {
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error('unexpected call'));

            expect(await transactionStore.saveTransactions(accountId)).toBe(true);

            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
        });
        it('calls updateTransactions mutation', async () => {
            const updates = [{id: 1, version: 2, memo: 'notes'}];
            const tableModel = new TransactionTableModel([], accountStore, categoryStore);
            jest.spyOn(tableModel, 'changes', 'get').mockReturnValue({updates});
            jest.spyOn(tableModel, 'update');
            transactionStore['transactionsByAccountId'].set(accountId, tableModel);
            const updatedTransactions = [newTx()];
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: {transactions: updatedTransactions}});

            expect(await transactionStore.saveTransactions(accountId)).toBe(true);

            expect(agent.graphql).toBeCalledWith(updateTxMutation, {accountId, updates});
            expect(tableModel.update).toBeCalledWith(updatedTransactions);
            expect(messageStore.addProgressMessage).toBeCalledWith(savingTransactions);
            expect(messageStore.removeProgressMessage).toBeCalledWith(savingTransactions);
        });
    });
});
