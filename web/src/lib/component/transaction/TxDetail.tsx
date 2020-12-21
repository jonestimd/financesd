import React from 'react';
import {RootStoreContext} from '../../store/RootStore';
import {ITransactionDetail} from '../../model/TransactionModel';
import {Currency, Shares} from '../../formats';
import {observer} from 'mobx-react-lite';

interface IProps {
    detail: ITransactionDetail;
}

const Category: React.FC<IProps> = observer(({detail: {relatedDetail, transactionCategoryId}}) => {
    const {accountStore, categoryStore} = React.useContext(RootStoreContext);
    if (relatedDetail) {
        return <>
            <i className='material-icons md-18'>send</i>
            <span>{accountStore.getAccount(relatedDetail.transaction.accountId).name}</span>
        </>;
    }
    if (transactionCategoryId) {
        return <>
            <i className='material-icons md-18'>payment</i>
            <span>{categoryStore.getCategory(transactionCategoryId).displayName}</span>
        </>;
    }
    return null;
});

// TODO group and memo

const TxDetail: React.FC<IProps> = ({detail}) => {
    return (
        <div className='detail chip' id={detail.id}>
            <Category detail={detail} />
            <span><Currency>{detail.amount}</Currency></span>
            {detail.assetQuantity ? <span>(<Shares>{detail.assetQuantity}</Shares>)</span> : null}
        </div>
    );
};

export default TxDetail;