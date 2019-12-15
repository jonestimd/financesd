import React from 'react';
import {MessageStore} from './MessageStore';
import {AccountStore} from './AccountStore';
import {CategoryStore} from './CategoryStore';
import {PayeeStore} from './PayeeStore';
import {TransactionStore} from './TransactionStore';
import {SecurityStore} from './SecurityStore';

declare global {
    // tslint:disable-next-line: interface-name
    interface Window {
        rootStore: RootStore;
    }
}

export class RootStore {
    messageStore = new MessageStore();
    accountStore = new AccountStore(this.messageStore);
    categoryStore = new CategoryStore(this.messageStore);
    payeeStore = new PayeeStore(this.messageStore);
    securityStore = new SecurityStore(this.messageStore);
    transactionStore = new TransactionStore(this);
}

export const RootStoreContext = React.createContext(new RootStore());
