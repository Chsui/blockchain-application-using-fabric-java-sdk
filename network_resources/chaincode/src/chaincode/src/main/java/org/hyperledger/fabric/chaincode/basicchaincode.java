package org.hyperledger.fabric.chaincode;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.google.protobuf.ByteString;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.hyperledger.fabric.chaincode.models.Ticket;
import org.hyperledger.fabric.shim.ChaincodeBase;
import org.hyperledger.fabric.shim.ChaincodeStub;

import java.util.List;

import static java.nio.charset.StandardCharsets.UTF_8;

public class basicchaincode extends ChaincodeBase {
    private static Log _logger = LogFactory.getLog(basicchaincode.class);
    private GsonBuilder builder = new GsonBuilder();
    private Gson gson = builder.create();

    @Override
    public Response init(ChaincodeStub stub) {
        try {
            List<String> args = stub.getStringArgs();
            if (args.size() != 5) {
                newErrorResponse("Error Incorrent arguments.");
            }
            Ticket ticket = new Ticket(args.get(0), args.get(1), args.get(2), args.get(3), args.get(4));
            String ticket_id = ticket.getId();
            stub.putState(ticket_id, gson.toJson(ticket).getBytes());
            return newSuccessResponse();
        } catch (Throwable e) {
            return newErrorResponse("Failed to create asset.");
        }
    }

    @Override
    public Response invoke(ChaincodeStub stub) {
        try {
            String func = stub.getFunction();
            List<String> params = stub.getParameters();
            if(func.equals("createTicket")) {
                return createTicket(stub, params);
            }
            if(func.equals("changeOwner")) {
                return changeOwner(stub, params);
            }
            if(func.equals("set")) {
                return set(stub, params);
            }
            if(func.equals("transfer")) {
                return transfer(stub, params);
            }
            if(func.equals("get")) {
                return get(stub, params);
            }
            return newErrorResponse("Invalid invoke function name.");
        } catch(Throwable e) {
            return newErrorResponse(e.getMessage());
        }
    }

    private Response createTicket(ChaincodeStub stub, List<String> args) {
        if(args.size() != 5) {
            throw new RuntimeException("Error Incorrect arguments");
        }
        Ticket ticket = new Ticket(args.get(0), args.get(1), args.get(2), args.get(3), args.get(4));
        String ticket_id = ticket.getId();
        try {
            stub.putState(ticket_id, gson.toJson(ticket).getBytes());
            return newSuccessResponse("Ticket Create (" + ticket_id + ")");
        } catch(Throwable e) {
            return newErrorResponse(e.getMessage());
        }
    }

    private Response changeOwner(ChaincodeStub stub, List<String> args) {
        if(args.size() != 2) {
            throw new RuntimeException("Error Incorrect arguments");
        }
        Ticket ticket = gson.fromJson(stub.getStringState(args.get(0)), Ticket.class);
        ticket.setOwner(args.get(1));
        String ticket_id = ticket.getId();
        try {
            stub.putState(ticket_id, gson.toJson(ticket).getBytes());
            return newSuccessResponse("Change Owner (" + ticket_id + ")");
        } catch(Throwable e) {
            return newErrorResponse(e.getMessage());
        }
    }

    private Response set(ChaincodeStub stub, List<String> args) {
        if(args.size() != 2) {
            throw new RuntimeException("Error Incorrect arguments");
        }
        stub.putStringState(args.get(0), args.get(1));
        return newSuccessResponse(args.get(1));
    }

    private Response get(ChaincodeStub stub, List<String> args) {
        if(args.size() != 1) {
            throw new RuntimeException("Incorrect number of arguments");
        }
        String key = args.get(0);
        String val = stub.getStringState(key);
        if(val == null) {
            return newErrorResponse(String.format("Error: state for %s is null", key));
        }
        return newSuccessResponse(val, ByteString.copyFrom(val, UTF_8).toByteArray());
    }

    private Response transfer(ChaincodeStub stub, List<String> args) {
        if(args.size() != 3) {
            throw new RuntimeException("Incorrect number of arguments");
        }
        String accountFromKey = args.get(0);
        String accountToKey = args.get(1);
        String accountFromValueStr = stub.getStringState(accountFromKey);
        if(accountFromValueStr == null) {
            return newErrorResponse(String.format("Entity %s not found", accountFromKey));
        }
        int accountFromValue = Integer.parseInt(accountFromValueStr);
        String accountToValueStr = stub.getStringState(accountFromValueStr);
        if(accountToValueStr == null) {
            return newErrorResponse(String.format("Entity %s not found", accountToKey));
        }
        int accountToValue = Integer.parseInt(accountToValueStr);
        int amount = Integer.parseInt(args.get(2));
        if(amount > accountFromValue) {
            return newErrorResponse(String.format("not enough money in account %s", accountFromKey));
        }
        accountFromValue -= amount;
        accountToValue += amount;
        _logger.info(String.format("new value of A: %s", accountFromValue));
        _logger.info(String.format("new value of B: %s", accountToValue));
        stub.putStringState(accountFromKey, Integer.toString(accountFromValue));
        stub.putStringState(accountToKey, Integer.toString(accountToValue));
        _logger.info("Transfer complete");
        return newSuccessResponse("invoke finished successfully", ByteString.copyFrom(accountFromKey + ":" + accountFromValue + " " + accountToKey + ": " + accountToValue, UTF_8).toByteArray());
    }

    public static void main(String[] args) {
        new basicchaincode().start(args);
    }

}
