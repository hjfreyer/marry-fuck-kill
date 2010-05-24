package com.mfk.web.client;

import com.mfk.web.shared.FieldVerifier;
import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.core.client.GWT;
import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.event.dom.client.KeyCodes;
import com.google.gwt.event.dom.client.KeyUpEvent;
import com.google.gwt.event.dom.client.KeyUpHandler;
import com.google.gwt.user.client.rpc.AsyncCallback;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.Label;
import com.google.gwt.user.client.ui.RootPanel;
import com.google.gwt.user.client.ui.TextBox;
import com.google.gwt.user.client.ui.VerticalPanel;

/**
 * Entry point classes define <code>onModuleLoad()</code>.
 */
public class MfkWeb implements EntryPoint {
	public static final String[] mfkText = { "?", "Marry", "Fuck", "Kill" };

	public enum Mfk {
		NONE, MARRY, FUCK, KILL
	};

	public static Mfk votes[] = {Mfk.NONE, Mfk.NONE, Mfk.NONE};
	
	private static Label statusLabel = new Label();
	public static void setStatus(String status) {
		MfkWeb.statusLabel.setText(status);
	}
	
	private static void addVoteButtons(String id, int itemNum) {
		Button marryButton = new Button("Marry");
		Button fuckButton = new Button("Fuck");
		Button killButton = new Button("Kill");
		marryButton.setWidth("200px");
		fuckButton.setWidth("200px");
		killButton.setWidth("200px");

		RootPanel rp = RootPanel.get(id);
		VerticalPanel vp = new VerticalPanel();
		rp.add(vp);
		
		vp.add(marryButton);
		vp.add(fuckButton);
		vp.add(killButton);

		marryButton.addClickHandler(VoteHandler.getHandler(Mfk.MARRY, itemNum));
		fuckButton.addClickHandler(VoteHandler.getHandler(Mfk.FUCK, itemNum));
		killButton.addClickHandler(VoteHandler.getHandler(Mfk.KILL, itemNum));
	}
	
	public void onModuleLoad() {
		// TODO(mkelly) fix the numbering scheme here
		MfkWeb.addVoteButtons("item1Vote", 0);
		MfkWeb.addVoteButtons("item2Vote", 1);
		MfkWeb.addVoteButtons("item3Vote", 2);
		
		final Button goButton = new Button("Go!");
		
		goButton.addClickHandler(new ClickHandler() {
			public void onClick(ClickEvent event) {
				MfkWeb.setStatus("Vote: 1: " + MfkWeb.votes[0] +
								  ", 2: " + MfkWeb.votes[1] +
								  ", 3: " + MfkWeb.votes[2]);
			}
		});
		
		RootPanel.get("status").add(MfkWeb.statusLabel);
		RootPanel.get("control").add(goButton);
		
		MfkWeb.setStatus("Welcome to MFK!");
	}
}

/**
 * Sets text on a given label.
 */
class VoteHandler implements ClickHandler {
	private String msg;
	private Label label;
	
	private MfkWeb.Mfk mfk;
	private int itemNum;

	// private constructor
	private VoteHandler() {
		super();
	}

	@Override
	public void onClick(ClickEvent event) {
		MfkWeb.votes[itemNum] = this.mfk;
		System.out.println("Vote: " + this.mfk + " item " + itemNum);
		MfkWeb.setStatus(this.mfk + " " + itemNum);
	}

	public static VoteHandler getHandler(MfkWeb.Mfk mfk, int itemNum) {
		VoteHandler h = new VoteHandler();
		h.mfk = mfk;
		h.itemNum = itemNum;
		return h;
	}
}