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
        return <span className='transfer'>
            <i className='material-icons md-18'>forward</i>
            {accountStore.getAccount(relatedDetail.transaction.accountId).name}
        </span>;
    }
    if (transactionCategoryId) {
        return <span className='category'>{categoryStore.getCategory(transactionCategoryId).displayName}</span>;
    }
    return null;
});

const Group: React.FC<{id?: number}> = observer(({id}) => {
    const {groupStore} = React.useContext(RootStoreContext);
    return typeof id === 'number' ? <span className='group'>{groupStore.getGroup(id).name}</span> : null;
});

const TxDetail: React.FC<IProps> = ({detail}) => {
    return (
        <div className='detail chip' id={detail.id}>
            <span><Currency>{detail.amount}</Currency></span>
            <Category detail={detail} />
            <Group id={detail.transactionGroupId} />
            {detail.assetQuantity ? <span className='shares'><Shares>{detail.assetQuantity}</Shares></span> : null}
            {detail.memo ? <span className='memo'>{detail.memo}</span> : null}
        </div>
    );
};

export default TxDetail;
