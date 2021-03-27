import React from 'react';
import {ITransactionDetail} from '../../model/TransactionModel';
import {Currency, Shares} from '../../formats';
import Category from './Category';
import Group from './Group';

interface IProps {
    detail: ITransactionDetail;
}

const TxDetail: React.FC<IProps> = ({detail}) => {
    return (
        <div className='detail chip' id={`${detail.id}`}>
            <span><Currency>{detail.amount}</Currency></span>
            <Category detail={detail} />
            <Group id={detail.transactionGroupId} />
            {detail.assetQuantity ? <span className='shares'><Shares>{detail.assetQuantity}</Shares></span> : null}
            {detail.memo ? <span className='memo'>{detail.memo}</span> : null}
        </div>
    );
};

export default TxDetail;
