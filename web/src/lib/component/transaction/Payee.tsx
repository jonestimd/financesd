import React from 'react';
import TransactionModel from 'src/lib/model/TransactionModel';
import {RootStoreContext} from '../../store/RootStore';

const Payee: React.FC<{transaction: TransactionModel}> = ({transaction: {payeeId}}) => {
    if (payeeId) {
        const {payeeStore} = React.useContext(RootStoreContext);
        return <span className='payee chip'><i className='material-icons md-18'>person</i>{payeeStore.getPayee(payeeId).name}</span>;
    }
    return null;
};

export default Payee;