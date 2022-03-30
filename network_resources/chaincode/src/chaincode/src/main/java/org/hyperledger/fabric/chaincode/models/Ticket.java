package org.hyperledger.fabric.chaincode.models;

public class Ticket {
    private String id;
    private String owner;
    private String showName;
    private String showDate;
    private String showLoc;
    private String position;
    private String price;

    public Ticket (String showName, String showDate, String showLoc, String position, String price) {
        // Ticker ID 에 대한 값을 어떠한 값으로 줄지 생각해볼 필요가 있음.
        this.id = "TickerID=" + showName + showDate + showLoc + position + price;
        this.owner = null;
        this.showName = showName;
        this.showDate = showDate;
        this.showLoc = showLoc;
        this.position = position;
        this.price = price;
    }

    public String getId() {
        return id;
    }

    public String getOwner() {
        return owner;
    }

    public String getShowName() {
        return showName;
    }

    public String getShowDate() {
        return showDate;
    }

    public String getShowLoc() {
        return showLoc;
    }
    public String getPosition() {
        return position;
    }

    public String getPrice() {
        return price;
    }

    public void setId(String id) {
        this.id = id;
    }

    public void setOwner(String owner) {
        this.owner = owner;
    }

    public void setShowName(String showName) {
        this.showName = showName;
    }

    public void setShowDate(String showDate) {
        this.showDate = showDate;
    }

    public void setShowLoc(String showLoc) {
        this.showLoc = showLoc;
    }

    public void setPosition(String position) {
        this.position = position;
    }

    public void setPrice(String price) {
        this.price = price;
    }

}
