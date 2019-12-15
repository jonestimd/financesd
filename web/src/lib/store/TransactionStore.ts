import agent from 'superagent';
import {TransactionModel, ITransaction} from '../model/TransactionModel';
import {IMessageStore} from './MessageStore';
import {flow, observable} from 'mobx';

const query = `query($accountId: ID) {
    transactions(accountId: $accountId) {
        date referenceNumber payeeId memo cleared
        details {
            id transactionCategoryId transactionGroupId memo amount assetQuantity
            relatedDetail {transaction {id accountId}}
        }
    }
}`;

interface ITransactionsResponse {
    body: {data: {transactions: ITransaction[]}};
}

const loadingTransactions = 'Loading transactions...';

function toModels(transactions: ITransaction[]): TransactionModel[] {
    const models = transactions.map(tx => new TransactionModel(tx)).sort(TransactionModel.compare);
    models.reduce((balance, model) => model.balance = balance + model.subtotal, 0);
    return models;
}

export class TransactionStore {
    private pendingAccounts: string[] = [];
    @observable
    private transactionsByAccountId: {[accountId: string]: TransactionModel[]} = {};
    private messageStore: IMessageStore;

    constructor(messageStore: IMessageStore) {
        this.messageStore = messageStore;
    }

    getTransactions(accountId: string): TransactionModel[] {
        return this.transactionsByAccountId[accountId] || [];
    }

    loadTransactions(accountId: string): void {
        if (!this.transactionsByAccountId[accountId] && this.pendingAccounts.indexOf(accountId) < 0) {
            this.messageStore.addProgressMessage(loadingTransactions);
            this._loadTransactions(accountId);
        }
    }

    private _loadTransactions = flow(function*(accountId: string) {
        this.pendingAccounts.push(accountId);
        try {
            const variables = {accountId};
            const {body: {data}}: ITransactionsResponse = yield agent.post('/finances/api/v1/graphql').send({query, variables});
            this.transactionsByAccountId[accountId] = toModels(data.transactions);
        } catch (err) {
            console.error('error gettting transactions', err);
        } finally {
            this.pendingAccounts.splice(this.pendingAccounts.indexOf(accountId), 1);
            this.messageStore.removeProgressMessage(loadingTransactions);
        }
    });
}