import React from 'react';
import {MessageStore, IMessageStore} from './MessageStore';
import {AccountStore} from './AccountStore';
import {CategoryStore, ICategoryStore} from './CategoryStore';
import {PayeeStore} from './PayeeStore';
import {TransactionStore} from './TransactionStore';

declare global {
    // tslint:disable-next-line: interface-name
    interface Window {
        rootStore: RootStore;
    }
}

export interface IRootStore {
    messageStore: IMessageStore;
    categoryStore: ICategoryStore;
}

export class RootStore {
    messageStore = new MessageStore();
    accountStore = new AccountStore(this.messageStore);
    categoryStore = new CategoryStore(this.messageStore);
    payeeStore = new PayeeStore(this.messageStore);
    transactionStore = new TransactionStore(this);
}

export const RootStoreContext = React.createContext(new RootStore());
