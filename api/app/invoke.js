const { Gateway, Wallets, TxEventHandler, GatewayOptions, DefaultEventHandlerStrategies, TxEventHandlerFactory } = require('fabric-network');
const fs = require('fs');
const path = require("path")
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const util = require('util')

const helper = require('./helper')


const invokeTransaction = async (channelName, chaincodeName, fcn, args, username, org_name, transientData) => {
    try {
        console.log(util.format('\n============ invoke transaction on channel %s ============\n', channelName));

        // load the network configuration
        // const ccpPath =path.resolve(__dirname, '..', 'config', 'connection-org1.json');
        // const ccpJSON = fs.readFileSync(ccpPath, 'utf8')
        const ccp = await helper.getCCP(org_name) //JSON.parse(ccpJSON);

        // Create a new file system based wallet for managing identities.
        const walletPath = await helper.getWalletPath(org_name) //path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        let identity = await wallet.get(username);
        if (!identity) {
            console.log(`An identity for the user ${username} does not exist in the wallet, so registering user`);
            await helper.getRegisteredUser(username, org_name, true)
            identity = await wallet.get(username);
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        const connectOptions = {
            wallet, identity: username, discovery: { enabled: true, asLocalhost: true },
            eventHandlerOptions: {
                commitTimeout: 100,
                strategy: DefaultEventHandlerStrategies.NETWORK_SCOPE_ALLFORTX
            },
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, connectOptions);

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork(channelName);
        const contract = network.getContract(chaincodeName);

        let result
        let message;
        if (fcn == "createAsset") {
            console.log("=========createAsset=========")
            console.log(fcn)
            result = await contract.submitTransaction(fcn, args[0], args[1], args[2], args[3]);
            message = `Successfully added the asset with the key ${args[0]}`
        } 
        else if (fcn == "createDEC") {
            console.log("=========creatignDEC=========")
            result = await contract.submitTransaction(fcn, args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8]);
            message = `Successfully added the DEC with the key ${args[0]}`
        }
        else if (fcn == "updateDEC") {
            console.log("=========updatingDEC=========")
            result = await contract.submitTransaction(fcn, args[0], args[1]);
        }   
        else if (fcn == "updateAssetStatus") {
            console.log("=========updatingAsset=========")
            result = await contract.submitTransaction(fcn, args[0], args[1]);
        }
        else if (fcn == "createPrivateAsset") {
            console.log(`Transient data is : ${transientData}`)
            let assetData = JSON.parse(transientData)
            console.log(`asset data is : ${JSON.stringify(assetData)}`)
            let key = Object.keys(assetData)[0]
            const transientDataBuffer = {}
            transientDataBuffer[key] = Buffer.from(JSON.stringify(assetData.asset))
            result = await contract.createTransaction(fcn)
                .setTransient(transientDataBuffer)
                .submit()
            message = `Successfully submitted transient data`
        }
        else if (fcn == "createEnergyRecord") {
            console.log("=========createEnergyRecord=========")
            console.log(fcn)
            result = await contract.submitTransaction(fcn, args[0]);
            message = `Successfully added the Energy Record`
        } 
        else if (fcn == "addWeeklyEnergyData") {
            console.log("=========addWeeklyEnergyData=========")
            console.log(fcn)
            result = await contract.submitTransaction(fcn, args[0], args[1], args[2]);
            message = `Successfully added the weekly data`
        } 
        else if (fcn == "createPerformanceContract") {
            console.log("=========createPerformanceContract=========")
            console.log(fcn)
            result = await contract.submitTransaction(fcn, args[0]);
            message = `Successfully added the performance contract`
        } 
        else if (fcn == "createModel") {
            console.log("=========createModel========")
            console.log(fcn)
            result = await contract.submitTransaction(fcn, args[0]);
            message = `Successfully added the model`
        } 
        else if (fcn == "reportPrediction") {
            console.log("=========reportPrediction========")
            console.log(fcn)
            result = await contract.submitTransaction(fcn, args[0]);
            message = `Successfully added the prediction`
        }
        else if (fcn == "trustUpdate") {
            console.log("=========trustUpdate========")
            console.log(fcn)
            result = await contract.submitTransaction(fcn, args[0]);
            message = `Successfully updated trust score`
        }  
        else {
            return `Invocation function not found. Function is ${fcn}`
        }

        await gateway.disconnect();

        result = JSON.parse(result.toString());

        let response = {
            message: message,
            result
        }

        return response;


    } catch (error) {

        console.log(`Getting error: ${error}`)
        return error.message

    }
}

exports.invokeTransaction = invokeTransaction;