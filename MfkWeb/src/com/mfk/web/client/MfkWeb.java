package com.mfk.web.client;

//import com.mfk.web.shared.FieldVerifier;
import com.google.gwt.core.client.EntryPoint;
import com.google.gwt.core.client.GWT;
import com.google.gwt.dev.util.collect.HashMap;
import com.google.gwt.event.dom.client.ClickEvent;
import com.google.gwt.event.dom.client.ClickHandler;
import com.google.gwt.event.dom.client.KeyCodes;
import com.google.gwt.event.dom.client.KeyUpEvent;
import com.google.gwt.event.dom.client.KeyUpHandler;
import com.google.gwt.http.client.Request;
import com.google.gwt.http.client.RequestBuilder;
import com.google.gwt.http.client.RequestCallback;
import com.google.gwt.http.client.RequestException;
import com.google.gwt.http.client.Response;
import com.google.gwt.json.client.JSONObject;
import com.google.gwt.json.client.JSONParser;
import com.google.gwt.json.client.JSONString;
import com.google.gwt.json.client.JSONValue;
import com.google.gwt.user.client.rpc.AsyncCallback;
import com.google.gwt.user.client.ui.Button;
import com.google.gwt.user.client.ui.DialogBox;
import com.google.gwt.user.client.ui.HTML;
import com.google.gwt.user.client.ui.Label;
import com.google.gwt.user.client.ui.RootPanel;
import com.google.gwt.user.client.ui.TextBox;
import com.google.gwt.user.client.ui.VerticalPanel;
import com.google.gwt.user.client.ui.Widget;

/**
 * Entry point classes define <code>onModuleLoad()</code>.
 */
public class MfkWeb implements EntryPoint {
	public static final String[] mfkText = { "?", "Marry", "Fuck", "Kill" };
	public static final int MAIN_VOTE_GROUP = 0;

	public enum Mfk {
		NONE, MARRY, FUCK, KILL
	};
	
	public static HTML entities[] = {new HTML(), new HTML(), new HTML()};
	public static VoteGroupHandler groups[] = {null, null, null};
	
	public static Button voteButton = new Button("Vote!");
	public static Button skipButton = new Button("Skip!");
	
	private static Label statusLabel = new Label();
	public static void setStatus(String status) {
		MfkWeb.statusLabel.setText(status);
	}
	
	private static void addVoteButtons(String id, int groupNum) {
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

		VoteGroupHandler group = new VoteGroupHandler(groupNum);
		marryButton.addClickHandler(new VoteChangeHandler(Mfk.MARRY, group));
		fuckButton.addClickHandler(new VoteChangeHandler(Mfk.FUCK, group));
		killButton.addClickHandler(new VoteChangeHandler(Mfk.KILL, group));
	}
	
	public static DialogBox makeErrorDialog(HTML html) {
	    final DialogBox dialog = new DialogBox();
	    dialog.setText("Error!");
	    dialog.setAnimationEnabled(false);
	    final Button closeButton = new Button("Close");
	    closeButton.getElement().setId("closeButton");
	    VerticalPanel dialogVPanel = new VerticalPanel();
	    dialogVPanel.addStyleName("dialogVPanel");
	    dialogVPanel.add(html);
	    dialogVPanel.setHorizontalAlignment(VerticalPanel.ALIGN_RIGHT);
	    dialogVPanel.add(closeButton);
	    dialog.setWidget(dialogVPanel);
	    closeButton.addClickHandler(new ClickHandler() {
	      public void onClick(ClickEvent event) {
	        dialog.hide();
	      }
	    });
	    return dialog;
	}
	
	public void onModuleLoad() {
		// Is the index-vs-label numbering confusing here? Maybe. Maybe not.
		MfkWeb.addVoteButtons("e1Vote", 0);
		MfkWeb.addVoteButtons("e2Vote", 1);
		MfkWeb.addVoteButtons("e3Vote", 2);
		
		MfkWeb.voteButton.setEnabled(false);
		
		voteButton.addClickHandler(new ClickHandler() {
			public void onClick(ClickEvent event) {
				MfkWeb.setStatus("Vote: 1: " + MfkWeb.groups[0] +
								  ", 2: " + MfkWeb.groups[1] +
								  ", 3: " + MfkWeb.groups[2]);
			}
		});
		
		RootPanel.get("status").add(MfkWeb.statusLabel);
		RootPanel.get("control").add(voteButton);
		RootPanel.get("control").add(skipButton);
		
		RootPanel.get("e1Display").add(entities[0]);
		RootPanel.get("e2Display").add(entities[1]);
		RootPanel.get("e3Display").add(entities[2]);
		
		final ClickHandler voteHandler = new LoadTripleHandler();
		skipButton.addClickHandler(voteHandler);
		// get an initial item
		voteHandler.onClick(null);
		MfkWeb.setStatus("Welcome to MFK!");
	}
	
	public static void setEntities(String one, String two, String three) {
		entities[0].setHTML(one);
		entities[1].setHTML(two);
		entities[2].setHTML(three);
	}

	public static void checkVoteStatus() {
		// a vote must exist for all buttons
		if (MfkWeb.groups[0].vote == MfkWeb.Mfk.NONE
				|| MfkWeb.groups[1].vote == MfkWeb.Mfk.NONE
				|| MfkWeb.groups[2].vote == MfkWeb.Mfk.NONE) {
			MfkWeb.voteButton.setEnabled(false);
		}
		else {
			// christ on a cracker. yes, this is it, in a nutshell.
			boolean valid = MfkWeb.groups[0].vote != MfkWeb.groups[1].vote
					&& MfkWeb.groups[1].vote != MfkWeb.groups[2].vote
					&& MfkWeb.groups[2].vote != MfkWeb.groups[0].vote;
			MfkWeb.voteButton.setEnabled(valid);
		}
	}
}

class LoadTripleHandler implements ClickHandler {
	private DialogBox errorDialog;
	
	public LoadTripleHandler() {
		this.errorDialog = MfkWeb.makeErrorDialog(
				new HTML("Error parsing JSON reply from server."));
	}
	
	@Override
	public void onClick(ClickEvent event) {
		String url = "/rpc/vote/";
		RequestBuilder builder = new RequestBuilder(RequestBuilder.GET, url);
		System.out.println("requesting new triple from " + url);
		
		try {
			builder.sendRequest(null, new RequestCallback() {
				@Override
				public void onError(Request request, Throwable exception) {
					System.out.println("Error retrieving new triple!");
				}

				@Override
				public void onResponseReceived(Request request, Response response) {
					JSONObject json = JSONParser.parse(response.getText()).isObject();
					System.out.println("Got JSONObject: " + json);
					try {
						MfkWeb.setEntities(
								json.get("one").isObject().get("name")
										.isString().stringValue(),
								json.get("two").isObject().get("name")
										.isString().stringValue(),
								json.get("three").isObject().get("name")
										.isString().stringValue());
					}
					catch (NullPointerException e) {
						System.out.println("Error parsing json reply!");
						errorDialog.center();
					}
				}
				
			});
		} catch (RequestException e) {
			e.printStackTrace();
		}
	}
}

/**
 * Processes a change in current vote state (assignment, in server-side
 * parlance).
 */
class VoteChangeHandler implements ClickHandler {
	private String msg;
	private Label label;
	
	private MfkWeb.Mfk vote;
	
	private VoteGroupHandler group;

	public VoteChangeHandler(MfkWeb.Mfk mfk, VoteGroupHandler group) {
		this.vote = mfk;
		this.group = group;
	}

	@Override
	public void onClick(ClickEvent event) {
		this.group.vote = this.vote;
		System.out.println("Vote: " + this.vote + " " + this.group);
		MfkWeb.setStatus(this.vote + " " + this.group);
		this.group.onClick(event);
	}
}

class VoteGroupHandler implements ClickHandler {
	public int num;
	public MfkWeb.Mfk vote;
	
	public VoteGroupHandler(int groupNum) {
		this.num = groupNum;
		MfkWeb.groups[this.num] = this;
	}
	
	@Override
	public void onClick(ClickEvent event) {
		System.out.println("Group " + this.num + " says hi");
		MfkWeb.checkVoteStatus();
	}
	
	public String toString() {
		return "<Group " + this.num + ": " + this.vote + ">";
	}
}