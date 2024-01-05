const AWS = require('aws-sdk')
const RDS = new AWS.RDSDataService()

/**
    * creates everything necessary to record a transaction bundle in RDS database
    *
    * A PTB (Purs Transaction Bundle) is a bundle of database records that are created as part of a Purs Transaction.
    * A single Transaction may include a promotion, it may include a Payment, it may include a fedNowPayment, but it always includes at least one LedgerEntry
    * @param {Object} userPurchaseInformation
    * @param {string} userPurchaseInformation.payor the id of the entity paying
    * @param {string} userPurchaseInformation.payee the id of the entity getting paid
    * @param {string} userPurchaseInformation.payorBankAccountID the bank account id of the entity paying
    * @param {string} userPurchaseInformation.payeeBankAccountID the bank account id of the entity getting paid
    * @param {string} userPurchaseInformation.dev the id of the developer
    * @param {integer} userPurchaseInformation.amount the amount being paid
    * @param {integer} userPurchaseInformation.interactionType 0 is mobile
    * @param {integer} userPurchaseInformation.paymentMethod 0 is fedNow 1 is card
    *
    * @param {Object}  promotionInformation
    * @param {integer}  promotionInformation.promoAmount if there is a promotion put an amount here
    * @param {string} sqlTransactionID the id of the sql transaction
    * @returns an array with the ledger entry id and the payment id
    */
const executeStandardPTOperations = async (userPurchaseInformation, promotionInformation, sqlTransactionID) => {
    const {
        payor,
        payee,
        payorBankAccountID,
        payeeBankAccountID,
        dev,
        amount,
        interactionType,
        paymentMethod,
    } = userPurchaseInformation

    const { promoAmount } = promotionInformation

    const ledgeEntries = []
    const paymentID = generateRandomBinary(32)
    let ledgerEntryID = generateRandomBinary(32)

    const idObj = {
        primaryPaymentID: paymentID,
        customerLedgerEntryID: ledgerEntryID,
    }

    ledgeEntries.push(ledgerEntryID)

    const params = {
        database: process.env.DATABASE,
        secretArn: process.env.SECRET_ARN,
        resourceArn: process.env.CLUSTER_ARN,
        transactionId: sqlTransactionID,
    }

    params.parameters = [
        {
            name: 'payerId',
            value: {
                blobValue: Buffer.from(payor, 'hex'),
            },
        },
        {
            name: 'payeeId',
            value: {
                blobValue: Buffer.from(payee, 'hex'),
            },
        },
        {
            name: 'paymentAmount',
            value: {
                doubleValue: amount,
            },
        },
        {
            name: 'interactionTypeId',
            value: {
                doubleValue: interactionType,
            },
        },
        {
            name: 'paymentId',
            value: {
                blobValue: Buffer.from(paymentID, 'hex'),
            },
        },
        {
            name: 'datePaid',
            value: {
                stringValue: paymentMethod === 0 && amount > 0 ? null : new Date().toISOString().slice(0, 19).replace('T', ' '),
                isNull: paymentMethod === 0 && amount > 0,
            },
        },
        {
            name: 'ledgerId',
            value: {
                blobValue: Buffer.from(ledgerEntryID, 'hex'),
            },
        },
        {
            name: 'developerId',
            value: {
                blobValue: Buffer.from(dev, 'hex'),
            },
        },
        {
            name: 'paymentMethod',
            value: {
                doubleValue: paymentMethod,
            },
        },
        {
            name: 'paymentStatus',
            value: {
                // the status is set to "completed" if the payment is a card or if the amount is 0; otherwise the status is set to pending(ie. fedNow)
                doubleValue: paymentMethod !== 0 || amount === 0 ? 4 : 5,
            },
        },
    ]

    params.parameters = [...params.parameters]
    params.sql = insertPaymentSQL


    await RDS.executeStatement({ ...params }).promise()

    // this is an additional step for processing fedNow payments
    if (paymentMethod === 0 && amount > 0) {
        const fedNowPaymentID = generateRandomBinary(32)

        params.parameters = [...params.parameters]
        params.parameters.push(
            {
                name: 'fedNowPaymentId',
                value: {
                    blobValue: Buffer.from(fedNowPaymentID, 'hex'),
                },
            },
            {
                name: 'payerAccountId',
                value: {
                    stringValue: payorBankAccountID
                }
            },
            {
                name: 'payeeAccountId',
                value: {
                    stringValue: payeeBankAccountID
                }
            }
        )

        params.sql = insertFedNowPaymentSQL
        await RDS.executeStatement({ ...params }).promise()

        idObj.primaryFedNowPaymentID = fedNowPaymentID
    }

    params.sql = insertLedgerEntrySQL
    await RDS.executeStatement({ ...params }).promise()

    if (promoAmount > 0) {
        ledgerEntryID = generateRandomBinary(32)
        ledgeEntries.push(ledgerEntryID)
        params.parameters = [
            {
                name: 'payerId',
                value: {
                    blobValue: Buffer.from(dev, 'hex'),
                },
            },
            {
                name: 'payeeId',
                value: {
                    blobValue: Buffer.from(payee, 'hex'),
                },
            },
            {
                name: 'amount',
                value: {
                    doubleValue: promoAmount,
                },
            },
            {
                name: 'interactionTypeId',
                value: {
                    doubleValue: interactionType,
                },
            },
            {
                name: 'ledgerId',
                value: {
                    blobValue: Buffer.from(ledgerEntryID, 'hex'),
                },
            },
            {
                name: 'developerId',
                value: {
                    blobValue: Buffer.from(dev, 'hex'),
                },
            },
        ]
        params.sql = insertPromoLedgerEntrySQL
        await RDS.executeStatement({ ...params }).promise()

        idObj.promotionLedgerEntryID = ledgerEntryID
    }

    const pursTransactionID = generateRandomBinary(32)

    params.parameterSets = ledgeEntries.map((pursPayment) => [{
        name: 'transactionId',
        value: {
            blobValue: Buffer.from(pursTransactionID, 'hex'),
        },
    }, {
        name: 'ledgerId',
        value: {
            blobValue: Buffer.from(pursPayment, 'hex'),
        },
    }])

    delete params.parameters
    params.sql = insertTransaction

    await RDS.batchExecuteStatement({ ...params }).promise()

    idObj.pursTransactionID = pursTransactionID

    return idObj
}
