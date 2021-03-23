import React from 'react';
import MessageStore from './MessageStore';
import AccountStore from './AccountStore';
import CategoryStore from './CategoryStore';
import PayeeStore from './PayeeStore';
import TransactionStore from './TransactionStore';
import SecurityStore from './SecurityStore';
import GroupStore from './GroupStore';
import AlertStore from './AlertStore';

declare global {
    interface Window {
        rootStore: RootStore;
    }
}

export class RootStore {
    alertStore = new AlertStore();
    messageStore = new MessageStore();
    accountStore = new AccountStore(this.messageStore);
    categoryStore = new CategoryStore(this.messageStore);
    groupStore = new GroupStore(this.messageStore);
    payeeStore = new PayeeStore(this.messageStore);
    securityStore = new SecurityStore(this.messageStore);
    transactionStore = new TransactionStore(this);
}

export const RootStoreContext = React.createContext(new RootStore());
