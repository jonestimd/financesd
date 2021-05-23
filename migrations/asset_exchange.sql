alter table transaction_category add column asset_exchange char(1) not null default 'N';

update transaction_category set asset_exchange = 'Y'
where code in ('Buy', 'Sell', 'Reinvest', 'Shares In', 'Shares Out');