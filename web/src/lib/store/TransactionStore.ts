import agent from 'superagent';
import {TransactionModel, ITransaction} from '../model/TransactionModel';
import {flow, observable} from 'mobx';
import {ICategoryStore} from './CategoryStore';
import {RootStore} from './RootStore';

const query = `query($accountId: ID) {
    transactions(accountId: $accountId) {
        id date referenceNumber payeeId securityId memo cleared
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

function toModels(transactions: ITransaction[], categoryStore: ICategoryStore): TransactionModel[] {
    const models = transactions.map(tx => new TransactionModel(tx, categoryStore)).sort(TransactionModel.compare);
    models.reduce((previous, model) => {
        model.previous = previous;
        return model;
    });
    return models;
}

export class TransactionStore {
    private pendingAccounts: string[] = [];
    @observable
    private transactionsByAccountId: {[accountId: string]: TransactionModel[]} = {};
    private rootStore: RootStore;

    constructor(rootStore: RootStore) {
        this.rootStore = rootStore;
    }

    getTransactions(accountId: string): TransactionModel[] {
        return this.transactionsByAccountId[accountId] || [];
    }

    loadTransactions(accountId: string): void {
        if (!this.transactionsByAccountId[accountId] && this.pendingAccounts.indexOf(accountId) < 0) {
            this.rootStore.messageStore.addProgressMessage(loadingTransactions);
            this._loadTransactions(accountId);
        }
    }

    private _loadTransactions = flow(function*(accountId: string) {
        this.pendingAccounts.push(accountId);
        try {
            const variables = {accountId};
            const {body: {data}}: ITransactionsResponse = yield agent.post('/finances/api/v1/graphql').send({query, variables});
            this.transactionsByAccountId[accountId] = toModels(data.transactions, this.rootStore.categoryStore);
        } catch (err) {
            console.error('error gettting transactions', err);
        } finally {
            this.pendingAccounts.splice(this.pendingAccounts.indexOf(accountId), 1);
            this.rootStore.messageStore.removeProgressMessage(loadingTransactions);
        }
    });
}