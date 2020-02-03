import React from 'react';
import TransactionModel from '../../model/TransactionModel';
import {RootStoreContext} from '../../store/RootStore';

const Security: React.FC<{transaction: TransactionModel}> = ({transaction: {securityId}}) => {
    if (securityId) {
        const {securityStore} = React.useContext(RootStoreContext);
        const security = securityStore.getSecurity(securityId);
        return <span className='security chip'><i className='material-icons md-18'>trending_up</i>{security.name}</span>;
    }
    return null;
};

export default Security;